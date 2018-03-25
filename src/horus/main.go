package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "horus/utils"
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
    gdaxUrl              = "https://api.gdax.com"
    currenciesPath       = "/currencies"
    cachedCurrenciesPath = "/data/currencies.json"
    // for cache vs fresh request
    freshData            = "fresh_data"
    cachedData           = "cached_data"

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
// - Use cli args to pass duration.
// - Write tests.
// - Abstract calls to exchange(s).

func main() {
    // use a ticker for this long running interval process
    // TODO Make time in seconds a cli arg
    timeInSeconds := time.Duration(10)
    ticker := time.NewTicker(time.Second * timeInSeconds)
    go func() {
        go launchGDAXCurrencyCheck() // launch first..
        fmt.Printf("Horus has been initiated with a ticker on %d seconds...\n", int(timeInSeconds))
        for range ticker.C {
            fmt.Printf("Horus: Launching currency check. Current time: %s\n", time.Now().Format(time.RFC1123))
            go launchGDAXCurrencyCheck()
        }
    }()

    // prevent this process from ending.
    keepRunning := make(chan interface{})
    <-keepRunning
}

func launchGDAXCurrencyCheck() {
    // paralelize goroutines
    proc := make(chan Currencies, 2)

    // freshCurrencies vs cachedCurrencies
    go requestGdaxCurrencies(currenciesPath, proc)
    go getCachedCurrencies(utils.GetGoPath() + cachedCurrenciesPath, proc)

    var callCount int
    var cachedCurrencies, freshCurrencies Currencies
    for currencies := range proc {
        callCount++
        switch collectionType := currencies.CollectionType; collectionType {
        case cachedData:
            cachedCurrencies = currencies
        case freshData:
            freshCurrencies = currencies
        }
        // _ = json.NewEncoder(os.Stdout).Encode(currencies.Collection)
        if callCount == 2 {
            close(proc)
        }
    }

    newCurrency, found := checkIfNewCurrencyFound(&cachedCurrencies, &freshCurrencies)
    if found {
        sendTwilioMessage(&newCurrency)
        os.Exit(0)
    } else {
        fmt.Println("Horus: No new currencies found. Checking back later...")
        fmt.Println()
    }
}

func sendTwilioMessage(newCurrency *Currency) {
    // TODO: Implement this function and exit process entirely once it ocurs
    fmt.Printf("Horus: sending notification that currency: %+v was just added..", newCurrency)
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
        os.Exit(1)
    }
    return newCurrency, true
}

// TODO: write tests - and more so this method.
func findNewlyAddedCurrency(cachedCurrencies *Currencies, freshCurrencies *Currencies) (Currency, bool) {
    // O + n^2
    var freshCurrency, cachedCurrency Currency
    var isThisNew bool
    for i := 0; i < len(freshCurrencies.Collection); i++ {
        freshCurrency = freshCurrencies.Collection[i]
        isThisNew = true
        for j := 0; j < len(cachedCurrencies.Collection); j++ {
            cachedCurrency = cachedCurrencies.Collection[j]
            // if there is no match between fresh and cached here, it means that is the new currency.
            if freshCurrency.Id == cachedCurrency.Id {
                isThisNew = false
                break
            }
        }
        if isThisNew {
            return freshCurrency, true
        }
    }
    return Currency{}, false
}

func requestGdaxCurrencies(path string, proc chan<- Currencies) {
    fmt.Printf("Horus: Requesting %s...\n", path)
    res, getErr := http.Get(gdaxUrl + path)
    if getErr != nil {
        log.Fatal(getErr)
        os.Exit(1)
        // what to do..
    }
    defer res.Body.Close()

    body, _         := ioutil.ReadAll(res.Body)
    freshCurrencies := Currencies{CollectionType: freshData}
    unmarshalErr    := json.Unmarshal(body, &freshCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
        os.Exit(1)
    }

    proc<- freshCurrencies
}

func getCachedCurrencies(path string, proc chan<- Currencies) {
    // defer use of a db. Just check data currencies.json file.
    // TODO: figure out interfaces better so i can return err or nil. (?)
    fmt.Println("Horus: Reading cached currencies ...")
    content, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }
    // todo we probably do need a store to write time of caching.
    cachedCurrencies := Currencies{CollectionType: cachedData}
    unmarshalErr := json.Unmarshal(content, &cachedCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
        os.Exit(1)
    }
    proc<- cachedCurrencies
}
