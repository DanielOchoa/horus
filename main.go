package main

import (
	"encoding/json"
	"fmt"
	"github.com/DanielOchoa/horus/config"
	"github.com/DanielOchoa/horus/http"
	"github.com/DanielOchoa/horus/twilio"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
)

// Part One.
// Horus

// A simple loop that checks gdax's available currencies. If a new currency is added - it needs to:
// 1. On request - check against a historical copy. If no copy is available:
//   1. Copy current copy store it somewhere (filesystem? memory?).
// 2. Check historical list against new request. If a new currency has been added:
//   1. Send text alert of currency to list of available phones.

// we only need id and/or name to check.
// walk through each freshCurrency and check against matching pair in cached currencies. If we can't find a
// match, we have a new currency in gdax.

const (
	gdaxUrl        = "https://api.pro.coinbase.com/"
	horusUA        = "HorusApp/1.0"
	currenciesPath = "/currencies"
	freshData      = "fresh_data"
	cachedData     = "cached_data"
)

//
// Custom Types
//

type Currency struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	MinSize string `json:"min_size"`
	Message string `json:"message"`
}

type Currencies struct {
	Collection     []Currency
	CollectionType string
}

// TODO:
// - Implement twilio messaging.
// - Logger util for prefixing stdout output with `Horus:`.
// - ~Use cli args to pass duration.~
// - Write tests.
// - Abstract http calls to exchange(s).
// - Move Currency types to it's own import.
// - Abstract flags into it's own flag composer(?).

// main here works as a never ending process. The only reason it ends is because:
// a) there was an error.
// b) We successfully found a new currency.
// The actual work is done inside a ticker interval anonymous function and the last part of main entails the
// mechanism to which maintains the process running indefinitely.

func main() {

	// .env file parsing into environment variables.
	envPath := config.GetFullProjectPath() + "/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal(err)
	}

	// setup cli flags. See func definition for flag options. Run `horus -h` to see options.
	tickerTime, cachedCurrenciesPath := config.SetupFlags()

	ticker := time.NewTicker(time.Second * time.Duration(tickerTime))

	go launchGDAXCurrencyCheck(cachedCurrenciesPath) // launch first..

	fmt.Printf("Horus: Horus has been initiated with a ticker of %d seconds...\n", tickerTime)

	for range ticker.C {
		fmt.Println("Horus:")
		fmt.Printf("Horus: Running currency check...\n")
		fmt.Printf("Horus: Current time: %s\n", time.Now().Format(time.RFC1123))
		fmt.Println("Horus:")
		go launchGDAXCurrencyCheck(cachedCurrenciesPath)
	}
}

func launchGDAXCurrencyCheck(cachedCurrenciesPath string) {
	// paralelize goroutines
	freshCurrencyCheckProc := make(chan Currencies, 1)
	cachedCurrencyFetchProc := make(chan Currencies, 1)

	parsedUrl, err := url.Parse(gdaxUrl)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{BaseURL: parsedUrl, UserAgent: horusUA}

	go requestGdaxCurrencies(client, currenciesPath, freshCurrencyCheckProc)
	go getCachedCurrencies(config.GetGoPath()+cachedCurrenciesPath, cachedCurrencyFetchProc)

	freshCurrencies := <-freshCurrencyCheckProc
	cachedCurrencies := <-cachedCurrencyFetchProc

	if newCurrency, found := checkIfNewCurrencyFound(&cachedCurrencies, &freshCurrencies); found {
		fmt.Printf("Horus: sending notification that currency: %+v was just added..\n", newCurrency)

		msg := fmt.Sprintf("Horus here.\nGDAX has posted a new coin!\n\nid: %s\nname: %s\n\nEnjoy!", newCurrency.Id, newCurrency.Name)
		toNumber := os.Getenv("TWILIO_TO_NUMBER")
		twilio.SendMessage(toNumber, msg)

		os.Exit(0)
	}

	fmt.Println("Horus: No new currencies found. Checking back later...")
}

func checkIfNewCurrencyFound(cachedCurrencies *Currencies, freshCurrencies *Currencies) (Currency, bool) {
	if len(cachedCurrencies.Collection) == len(freshCurrencies.Collection) {
		return Currency{}, false
	}
	fmt.Println()
	fmt.Println("Horus: ========================== OMG | WE GOT A NEW COIN! | OMG ===========================")
	fmt.Println()
	// figure out which one is the new one..
	// TODO: We'll have to write tests for this..
	// Note that if more than 1 currency is added.. it doens't matter for now.
	newCurrency, found := findNewlyAddedCurrency(cachedCurrencies, freshCurrencies)
	if !found {
		// we SHOULD have picked up a currency so obviously something went horribly wrong...
		log.Fatal("Yeah dawg the `findNewlyAddedCurrency function is not working properly...`")
	}
	return newCurrency, true
}

// Compares two different lists of Currencies, and finds the first that doesn't exist between the two lists.
//
// A cleaner implementation would be to convert the structs to maps and then substract key/value pairs (? maybe).
func findNewlyAddedCurrency(cachedCurrencies *Currencies, freshCurrencies *Currencies) (Currency, bool) {
	cachedCurrenciesMap := make(map[Currency]bool)
	for _, cachedCurrency := range cachedCurrencies.Collection {
		cachedCurrenciesMap[cachedCurrency] = true
	}
	for _, freshCurrency := range freshCurrencies.Collection {
		if _, ok := cachedCurrenciesMap[freshCurrency]; !ok {
			return freshCurrency, true
		}
	}
	return Currency{}, false
}

func requestGdaxCurrencies(c http.Client, path string, proc chan<- Currencies) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		log.Fatal(err)
	}
	freshCurrencies := Currencies{CollectionType: freshData}
	res, err := c.Do(req, &freshCurrencies.Collection)
	if err != nil {
		fmt.Printf("Horus: WARNING - `Do` request error: %q\n", err)
	}
	defer res.Body.Close()

	proc <- freshCurrencies
}

func getCachedCurrencies(path string, proc chan<- Currencies) {
	// defer use of a db. Just check data currencies.json file.
	// TODO: figure out interfaces better so i can return err or nil. (?)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	// todo we probably do need a store to write time of caching.
	cachedCurrencies := Currencies{CollectionType: cachedData}
	unmarshalErr := json.Unmarshal(content, &cachedCurrencies.Collection)
	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}
	proc <- cachedCurrencies
}
