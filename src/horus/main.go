package main

import (
    "net/http"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
)
// Part One.
// Horus

// A simple loop that checks gdax's available currencies. If a new currency is added - it needs to:
// 1. On request - check against a historical copy. If no copy is available:
//   1. Copy current copy store it somewhere (filesystem? memory?).
// 2. Check historical list against new request. If a new currency has been added:
//   1. Send text alert of currency to list of available phones.

const (
    gdaxUrl = "https://api.gdax.com"
    currenciesPath = "/currencies"
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
    requestGdaxCurrencies(currenciesPath)
}

func requestGdaxCurrencies(path string) {
    fmt.Printf("Horus: Requesting %s...\n", path)
    res, getErr := http.Get(gdaxUrl + path)
    if getErr != nil {
        log.Fatal(getErr)
        // what to do..
    }
    defer res.Body.Close()

    body, _      := ioutil.ReadAll(res.Body)
    currencies   := Currencies{}
    unmarshalErr := json.Unmarshal(body, &currencies.Collection)
    if unmarshalErr != nil {
        log.Fatal(unmarshalErr)
    }

    fmt.Println("Horus: Parsed json: ")
    for i, curr := range currencies.Collection {
        if i == 0 {
            fmt.Println("[")
        }
        fmt.Print("  ")
        _ = json.NewEncoder(os.Stdout).Encode(curr)
        if i + 1 == len(currencies.Collection) {
            fmt.Println("]")
        }
    }

    //_ = _stdoutJSON(currencies.Collection)
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
