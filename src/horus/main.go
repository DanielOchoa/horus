package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "time"
    // local
    "horus/utils"
)
// Part One.
// Horus

// A simple loop that checks gdax's available currencies. If a new currency is added - it needs to:
// 1. On request - check against a historical copy. If no copy is available:
//   1. Copy current copy store it somewhere (filesystem? memory?).
// 2. Check historical list against new request. If a new currency has been added:
//   1. Send text alert of currency to list of available phones.

const (
    gdaxUrl              = "https://api.gdax.com"
    currenciesPath       = "/currencies"
    cachedCurrenciesPath = "/data/currencies.json"
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
    Collection  []Currency
    CollectedOn time.Time
}

// TODO:
// storage device..
// check if payload is already saved.
// no? save it. EOF
// yes? check saved payload with new payload.
// get difference.
// parse differences by currency object.
// setup messenger lib
// send msg with new currencies to stored emails (db!).
// update cached version!

func main() {
    // paralelize goroutines
    proc := make(chan Currencies, 2)

    // freshCurrencies vs cachedCurrencies
    go requestGdaxCurrencies(currenciesPath, proc)
    go getCachedCurrencies(utils.GetGoPath() + cachedCurrenciesPath, proc)

    callCount := 0
    for res := range proc {
        callCount++
        fmt.Printf("%q", res.Collection)
        if callCount == 2 {
            close(proc)
        }
    }

    // _ = json.NewEncoder(os.Stdout).Encode(curr) // <- sample stdout of json
    // we only need id and/or name to check.
    // walk through each freshCurrency and check against matching pair in cached currencies. If we can't find a
    // match, we have a new currency in gdax.

}

func requestGdaxCurrencies(path string, proc chan<- Currencies) {
    fmt.Printf("Horus: Requesting %s...\n", path)
    res, getErr := http.Get(gdaxUrl + path)
    if getErr != nil {
        log.Fatal(getErr)
        // what to do..
    }
    defer res.Body.Close()

    body, _         := ioutil.ReadAll(res.Body)
    freshCurrencies := Currencies{CollectedOn: time.Now()}
    unmarshalErr    := json.Unmarshal(body, &freshCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
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
    cachedCurrencies := Currencies{}
    unmarshalErr := json.Unmarshal(content, &cachedCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
    }
    proc<- cachedCurrencies
}
