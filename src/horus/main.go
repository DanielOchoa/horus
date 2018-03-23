package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
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
    Collection []Currency
}


func main() {
    // freshCurrencies vs cachedCurrencies
    _ = requestGdaxCurrencies(currenciesPath)
    gopath := utils.GetGoPath()
    cachedCurrencies := getCachedCurrencies(gopath + cachedCurrenciesPath)
    fmt.Println("Now showing cached ones...")
    json.NewEncoder(os.Stdout).Encode(cachedCurrencies.Collection)
    // we only need id and/or name to check.
    // walk through each freshCurrency and check against matching pair in cached currencies. If we can't find a
    // match, we have a new currency in gdax.

}

func requestGdaxCurrencies(path string) Currencies {
    fmt.Printf("Horus: Requesting %s...\n", path)
    res, getErr := http.Get(gdaxUrl + path)
    if getErr != nil {
        log.Fatal(getErr)
        // what to do..
    }
    defer res.Body.Close()

    body, _         := ioutil.ReadAll(res.Body)
    freshCurrencies := Currencies{}
    unmarshalErr    := json.Unmarshal(body, &freshCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
    }

    fmt.Println("Horus: Parsed json: ")
    for i, curr := range freshCurrencies.Collection {
        // simple stdout json formatting
        if i == 0 {
            fmt.Println("[")
        }
        fmt.Print("  ")
        _ = json.NewEncoder(os.Stdout).Encode(curr)
        if i + 1 == len(freshCurrencies.Collection) {
            fmt.Println("]")
        }
    }

    return freshCurrencies



    // TODO:
    // storage device..
    // check if payload is already saved.
    // no? save it. EOF
    // yes? check saved payload with new payload.
    // get difference.
    // parse differences by currency object.
    // setup messenger lib
    // send msg with new currencies to stored emails (db!).
}

func getCachedCurrencies(path string) Currencies {
    // defer use of a db. Just check data currencies.json file.
    // TODO: figure out interfaces better so i can return err or nil. (?)
    content, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }
    cachedCurrencies := Currencies{}
    unmarshalErr := json.Unmarshal(content, &cachedCurrencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
    }
    return cachedCurrencies
}
