# HORUS

Horus is a platform to automate various tasks for crypto day trading.

## Step 1

Build a script that pings GDAX currencies endpoint to see if a new cryptocurrency has been added to their platform.
If so, notify me by msn when this occurs. The script will run this check every 10 minutes by default.

More features coming soon.

### Running

You'll need go installed.

`go run main.go -time 10`

#### Script arguments

 - The `time` argument is optional and it specifies the interval time in seconds between every currency endpoint
     check vs a local cached copy of it.
 - The `cachedCurrenciesPath` is optional. It defaults to the path of `/data/currencies.json`. Useful if you want to
     test with a modified json file in order to functionally test the script will work when needed.

### Running all tests

`go test ./...`. Note that while it runs all tests, it also runs dependencies test (or all packages in your `go/src` location for that matter).

## Author

Daniel Ochoa
