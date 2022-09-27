package main

import (
	"encoding/json"
	"fmt"
	loggly "github.com/JamesPEarly/loggly"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Symbol struct {
	Name                 string          `json:"Name"`
	StockSymbol          string          `json:"StockSymbol"`
	Price                int             `json:"Price"`
	DollarChange         int             `json:"DollarChange"`
	PercentChange        int             `json:"PercentChange"`
	PreviousClose        int             `json:"PreviousClose"`
	Open                 int             `json:"Open"`
	BidPrice             int             `json:"BidPrice"`
	BidQuantity          int             `json:"BidQuantity"`
	AskPrice             int             `json:"AskPrice"`
	AskQuantity          int             `json:"AskQuantity"`
	DayRangeLow          int             `json:"DayRangeLow"`
	DayRangeHigh         int             `json:"DayRangeHigh"`
	YearRangeLow         int             `json:"YearRangeLow"`
	YearRangeHigh        int             `json:"YearRangeHigh"`
	Volume               int             `json:"Volume"`
	AverageVolume        int             `json:"AverageVolume"`
	MarketCap            int             `json:"MarketCap"`
	Beta                 int             `json:"Beta"`
	PriceEarningsRatio   int             `json:"PriceEarningsRatio"`
	EarningsPerShare     int             `json:"EarningsPerShare"`
	EarningsDate         string          `json:"EarningsDate"`
	ForwardDividend      int             `json:"ForwardDividend"`
	ForwardDividendYield int             `json:"ForwardDividendYield"`
	ExDividendDate       int             `json:"ExDividendDate"`
	YearTargetEstimate   int             `json:"YearTargetEstimate"`
	QueriedSymbol        string          `json:"QueriedSymbol"`
	DataCollectedOn      time.Time       `json:"DataCollectedOn"`
	Stats                []SummaryStruct `json:"Summary"`
}

type SummaryStruct struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func main() {
	os.Setenv("LOGGLY_TOKEN", "e4a25bf2-e2cc-4771-95c8-b9a68c55bc11")
	client := loggly.New("anubhav")
	fmt.Println("Enter Ticker: ")
	var name string
	fmt.Scanln(&name)

	// Calling API
	req, err := http.NewRequest(
		http.MethodGet, "https://api.aletheiaapi.com/StockData?symbol="+name+"&summary=true&statistics=false",
		nil,
	)
	if err != nil {
		client.EchoSend("error", "Failed with error: "+err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("key", ("9765EE5F17A04F03B9A29C3DBBC698A3"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		client.EchoSend("error sending HTTP request: %v", err.Error())
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		client.EchoSend("error reading HTTP response body: %v", err.Error())
	}
	//	log.Println("We got the response:", string(responseBytes))

	// Parse the JSON and display info
	var symbol Symbol
	json.Unmarshal(responseBytes, &symbol)
	fmt.Println(string(responseBytes))

	// Send success message with response size
	var respSize string = strconv.Itoa(len(responseBytes))
	logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
	if logErr != nil {
		fmt.Println("err: ", logErr)
	}

}
