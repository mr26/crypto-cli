# Crypto CLI
Crypto CLI is a simple CLI tool used to retrieve cryptocurrency data and display it on your terminal.  All data is retrieved real-time 
and up-to-date from the coinmarketcap.com API's.

## Installation
`go get -u github.com/mr26/crypto-cli`


##Usage
```
crypto-cli -h
get-currency-data
  -name string
    	Name of the currency you want to see the current data for.
  -symbol string
    	Symbol of the currency you want to see data for (ex; BTC, ETH, DOT)

get-currency-symbol
  -name string
    	Name of the currency you want to get the symbol for.

get-market-data

```

##Build
Simply run `go build` and you will be good to go.  You can then move the crypto-cli binary to a directory specified in your environment's $PATH
in order to be able to execute it without having to specify a path.