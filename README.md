# HORUS

Horus eye watches all...

## Episode 1

 - Setup main to ping `gdax/currencies` api to see if a new cryptocurrency has been added and notify you by msn
     whenever this occurs.

More features coming soon.

### Running

You'll need go installed.

`go run main.go -time 600`

#### Script arguments

 - The `time` argument is optional and it specifies the interval time in seconds between api hits.
 - The `cachedCurrenciesPath` is optional. It defaults to the path of `/data/currencies.json`. Useful if you want to
test with a modified json file in order to functionally test the script works.

### Running all tests

`go test ./...`. Note that while it runs all tests, it also runs dependencies test (or all packages in your `go/src` location for that matter).

## Author

Daniel Ochoa
