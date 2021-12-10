package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type CMarketCapID struct {
	Id     int
	Name   string
	Symbol string
	Slug   string
}

type CMarketIdData struct {
	Data []CMarketCapID
}

type CMarketUsd struct {
	Price              float64
	Volume_24h         float64
	Volume_change_24h  float64
	Percent_change_1h  float64
	Percent_change_24h float64
	Percent_change_7d  float64
	Market_cap         float64
	Last_updated       string
}

type CMarketQuote struct {
	Usd CMarketUsd
}

type CMarketListing struct {
	Name               string
	Slug               string
	Id                 int
	Symbol             string
	Circulating_supply float64
	Cmc_rank           int
	Quote              CMarketQuote
}

type CMarketListings struct {
	Data []CMarketListing
}

// Set this environment variable locally.
var CMARKETCAP_API_KEY string = os.Getenv("CMARKETCAP_API_KEY")

func get_id(name_or_symbol string) (v_id int) {
	// returns coinmarketcap ID, this is needed in order to reliably search for data related to a currency.
	fmt.Println("Retrieving Id...")
	var cmarket CMarketIdData
	name_or_symbol = strings.ToLower(name_or_symbol)
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/map", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header = http.Header{
		"X-CMC_PRO_API_KEY": []string{CMARKETCAP_API_KEY},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server.")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(respBody)), &cmarket)
	for _, v := range cmarket.Data {
		if strings.ToLower(v.Name) == name_or_symbol || strings.ToLower(v.Slug) == name_or_symbol || strings.ToLower(v.Symbol) == name_or_symbol {
			return v.Id
		}
	}
	return
}

func get_currency_data(name_or_symbol string) (currency CMarketListing) {
	// returns currency data for a single currency in raw api format.
	var c_listings CMarketListings
	c_id := get_id(name_or_symbol)
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header = http.Header{
		"X-CMC_PRO_API_KEY": []string{CMARKETCAP_API_KEY},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server.")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(respBody)), &c_listings)
	for _, v := range c_listings.Data {
		if v.Id == c_id {
			fmt.Println(v)
			return v
		}
	}
	return
}

func display_currency_data(name_or_symbol string) {
	// displays currency data for a single currency that is provided by the user.
	currency := get_currency_data(name_or_symbol)
	zone, _ := time.Now().Zone()
	est, err1 := time.LoadLocation(zone)
	if err1 != nil {
		fmt.Println(err1)
	}
	tm, err2 := time.ParseInLocation(time.RFC3339, currency.Quote.Usd.Last_updated, est)
	if err2 != nil {
		fmt.Println(err2)
	}
	p := message.NewPrinter(language.English)
	data := [][]string{
		{fmt.Sprintf("%v", currency.Cmc_rank), currency.Name, currency.Symbol, p.Sprintf("$%.2f", currency.Quote.Usd.Price),
			fmt.Sprintf("%.2f", currency.Quote.Usd.Percent_change_1h), fmt.Sprintf("%.2f", currency.Quote.Usd.Percent_change_24h),
			fmt.Sprintf("%.2f", currency.Quote.Usd.Percent_change_7d), p.Sprintf("$%.2f", currency.Quote.Usd.Market_cap),
			p.Sprintf("%.f", currency.Circulating_supply), tm.In(est).String()},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"CMC Rank", "Name", "Symbol", "Price", "1h %", "24h %", "7d %", "Market Cap", "Circulating Supply", "Last Updated"})
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

func get_market_data() []CMarketListing {
	// Returns coinmarketdata in raw API format. Data is then parsed through in display_market_data function.
	var c_listings CMarketListings
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header = http.Header{
		"X-CMC_PRO_API_KEY": []string{CMARKETCAP_API_KEY},
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server.")
		os.Exit(1)
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(string(respBody)), &c_listings)
	return c_listings.Data
}

func display_market_data() {
	// Displays market data for top 100 coins in coinmarketcap.com.
	market_data := [][]string{}
	market_listings := get_market_data()
	zone, _ := time.Now().Zone()
	est, err := time.LoadLocation(zone)
	if err != nil {
		fmt.Println(err)
	}
	p := message.NewPrinter(language.English)
	for _, v := range market_listings {
		if v.Cmc_rank == 101 {
			break
		}
		tm, err := time.ParseInLocation(time.RFC3339, v.Quote.Usd.Last_updated, est)
		if err != nil {
			fmt.Println(err)
		}
		data := []string{
			fmt.Sprintf("%v", v.Cmc_rank), v.Name, v.Symbol, p.Sprintf("$%.2f", v.Quote.Usd.Price),
			fmt.Sprintf("%.2f", v.Quote.Usd.Percent_change_1h), fmt.Sprintf("%.2f", v.Quote.Usd.Percent_change_24h),
			fmt.Sprintf("%.2f", v.Quote.Usd.Percent_change_7d), p.Sprintf("$%.2f", v.Quote.Usd.Market_cap),
			p.Sprintf("%.f", v.Circulating_supply), tm.In(est).String()}
		market_data = append(market_data, data)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetRowSeparator("-")
	table.SetHeader([]string{"CMC Rank", "Name", "Symbol", "Price", "1h %", "24h %", "7d %", "Market Cap", "Circulating Supply", "Last Updated"})
	for _, v := range market_data {
		table.Append(v)
	}
	table.Render()
}

func get_currency_symbol(name string) {
	// Displays currency symbol based on the currency name provided by user.
	data := get_market_data()
	for _, v := range data {
		if strings.ToLower(v.Name) == name || strings.ToLower(v.Slug) == name {
			fmt.Println(v.Symbol)
		}
	}
}

func get_currency_data_func(get_currency_data_cmd *flag.FlagSet, currency_name *string, currency_symbol *string) {
	get_currency_data_cmd.Parse(os.Args[2:])
	if *currency_name != "" {
		display_currency_data(*currency_name)
	} else if *currency_symbol != "" {
		display_currency_data(*currency_symbol)
	} else {
		fmt.Println("Please provide either the --symbol or --name flag along with a corresponding currency symbol or name.")
		os.Exit(1)
	}
}

func get_currency_symbol_func(get_currency_symbol_cmd *flag.FlagSet, cur_name *string) {
	get_currency_symbol_cmd.Parse(os.Args[2:])
	if *cur_name == "" {
		fmt.Println("Please provide the --name flag with the name of the currency you want the symbol for (e.g. --name 'bitcoin').")
		os.Exit(1)
	} else {
		get_currency_symbol(*cur_name)
	}
}

func print_help_message(get_currency_data_cmd *flag.FlagSet, get_currency_symbol_cmd *flag.FlagSet) {
	fmt.Println("Expected 'get-currency-data', 'get-market-data', or 'get-currency-symbol' subcommands.")
	fmt.Println("get-currency-data")
	get_currency_data_cmd.PrintDefaults()
	fmt.Println("\nget-currency-symbol")
	get_currency_symbol_cmd.PrintDefaults()
	fmt.Println("\nget-market-data")
	fmt.Println("")
}

func main() {
	get_currency_data_cmd := flag.NewFlagSet("get-currency-data", flag.ExitOnError)
	currency_name := get_currency_data_cmd.String("name", "", "Name of the currency you want to see the current data for.")
	currency_symbol := get_currency_data_cmd.String("symbol", "", "Symbol of the currency you want to see data for (ex; BTC, ETH, DOT)")
	get_market_data_cmd := flag.NewFlagSet("get-market-data", flag.ExitOnError)
	get_currency_symbol_cmd := flag.NewFlagSet("get-currency-symbol", flag.ExitOnError)
	cur_name := get_currency_symbol_cmd.String("name", "", "Name of the currency you want to get the symbol for.")
	flag.Usage = func() {
		fmt.Println("get-currency-data")
		get_currency_data_cmd.PrintDefaults()
		fmt.Println("\nget-currency-symbol")
		get_currency_symbol_cmd.PrintDefaults()
		fmt.Println("\nget-market-data")
		fmt.Println("")
	}
	flag.Parse()
	if len(os.Args) < 2 {
		print_help_message(get_currency_data_cmd, get_currency_symbol_cmd)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "get-currency-data":
		get_currency_data_func(get_currency_data_cmd, currency_name, currency_symbol)
	case "get-market-data":
		get_market_data_cmd.Parse([]string{})
		display_market_data()
	case "get-currency-symbol":
		get_currency_symbol_func(get_currency_data_cmd, cur_name)
	default:
		fmt.Println("The subcommand you provided was invalid. Please provide a valid subcommand.")
		os.Exit(1)
	}
}
