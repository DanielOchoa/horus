package main

import "testing"

func TestFindNewlyAddedCurrency(t *testing.T) {
	btc := &Currency{
		Id:   "btc",
		Name: "Bitcoin",
	}
	eth := &Currency{
		Id:   "eth",
		Name: "Ethereum",
	}
	usd := &Currency{
		Id:   "usd",
		Name: "United States Dollars",
	}
	xrp := &Currency{
		Id:   "xrp",
		Name: "Ripple",
	}

	newCurrenciesList := &Currencies{Collection: []Currency{*btc, *usd, *eth, *xrp}}
	cachedCurrenciesList := &Currencies{Collection: []Currency{*btc, *usd, *eth}}

	newCurrency, found := findNewlyAddedCurrency(cachedCurrenciesList, newCurrenciesList)
	if !found {
		t.Error("Expected to have found a new currency")
	}
	if newCurrency.Id != xrp.Id {
		t.Error("The new currency found must have Id of", xrp.Id)
	}

	// now check that it should not find a new currency
	newCurrenciesList = &Currencies{Collection: []Currency{*btc, *usd, *eth}}
	cachedCurrenciesList = &Currencies{Collection: []Currency{*btc, *usd, *eth}}

	newCurrency, found = findNewlyAddedCurrency(newCurrenciesList, cachedCurrenciesList)
	if found {
		t.Error("Expect to have not found any currencies.")
	}
	if newCurrency.Id != "" {
		t.Error("Expect `newCurrency.Id` to be empty.")
	}
}
