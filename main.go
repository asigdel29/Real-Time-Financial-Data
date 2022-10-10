package main

import (
	"encoding/json"
	"fmt"
	"github.com/JamesPEarly/loggly"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Symbol struct {
	Name                 string    `json:"Name"`
	StockSymbol          string    `json:"StockSymbol"`
	Price                int       `json:"Price"`
	DollarChange         int       `json:"DollarChange"`
	PercentChange        int       `json:"PercentChange"`
	PreviousClose        int       `json:"PreviousClose"`
	Open                 int       `json:"Open"`
	BidPrice             int       `json:"BidPrice"`
	BidQuantity          int       `json:"BidQuantity"`
	AskPrice             int       `json:"AskPrice"`
	AskQuantity          int       `json:"AskQuantity"`
	DayRangeLow          int       `json:"DayRangeLow"`
	DayRangeHigh         int       `json:"DayRangeHigh"`
	YearRangeLow         int       `json:"YearRangeLow"`
	YearRangeHigh        int       `json:"YearRangeHigh"`
	Volume               int       `json:"Volume"`
	AverageVolume        int       `json:"AverageVolume"`
	MarketCap            int       `json:"MarketCap"`
	Beta                 int       `json:"Beta"`
	PriceEarningsRatio   int       `json:"PriceEarningsRatio"`
	EarningsPerShare     int       `json:"EarningsPerShare"`
	EarningsDate         string    `json:"EarningsDate"`
	ForwardDividend      int       `json:"ForwardDividend"`
	ForwardDividendYield int       `json:"ForwardDividendYield"`
	ExDividendDate       int       `json:"ExDividendDate"`
	YearTargetEstimate   int       `json:"YearTargetEstimate"`
	QueriedSymbol        string    `json:"QueriedSymbol"`
	DataCollectedOn      time.Time `json:"DataCollectedOn"`
}

type Symboldata struct {
	Symbols []Symbol `json:"Symbols"`
}

type Item struct {
	Time    time.Time
	Name    string
	Symbols []byte
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	os.Setenv("LOGGLY_TOKEN", goDotEnvVariable("LOGGLY_TOKEN"))
	os.Setenv("AWS_ACCESS_KEY_ID", goDotEnvVariable("AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", goDotEnvVariable("AWS_SECRET_ACCESS_KEY"))
	apikey := goDotEnvVariable("API_Key")

	for {
		stocks := [10]string{"TSLA", "AAPL", "MSFT", "GOOGL", "NIO", "NVDA", "MRNA", "NKLA", "FB", "AMD"}
		for _, element := range stocks {
			weekday := time.Now().Weekday()
			daycheck := (int(weekday))
			if daycheck == 7 || daycheck == 0 {
				fmt.Println("Market Closed")
			} else {
				var tag string = element
				client := loggly.New(tag)
				req, err := http.NewRequest(
					http.MethodGet, "https://api.aletheiaapi.com/StockData?symbol="+element+"&summary=true&statistics=false",
					nil,
				)

				if err != nil {
					client.EchoSend("error", "Failed with error: "+err.Error())
				}

				req.Header.Add("Accept", "application/json")
				req.Header.Add("key", (apikey))

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					client.EchoSend("error sending HTTP request: %v", err.Error())
				}

				responseBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					client.EchoSend("error reading HTTP response body: %v", err.Error())
				}
				//	log.Println("We got the response:", string(responseBytes))

				var symbol Symbol
				json.Unmarshal(responseBytes, &symbol)
				//formattedData, _ := json.MarshalIndent(symbol, "    ", "    ")
				//fmt.Println(formattedData)
				fmt.Println(string(responseBytes))

				var respSize string = strconv.Itoa(len(responseBytes))
				logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
				if logErr != nil {
					fmt.Println("err: ", logErr)
				}
				/*
					sess, err := session.NewSession(&aws.Config{
						Region: aws.String("us-east-1")},
					)
					if err != nil {
						log.Fatalf("Error initializing AWS: %s", err)
					}

					svc := dynamodb.New(sess)
					var item Item
					item.Time = symbol.DataCollectedOn
					item.Name = symbol.Name
					item.Symbols = responseBytes

					av, err := dynamodbattribute.MarshalMap(item)
					if err != nil {
						log.Fatalf("Error marshalling %s", err)
					}

					tableName := "Stock Summary"
					input := &dynamodb.PutItemInput{
						Item:      av,
						TableName: aws.String(tableName),
					}

					_, err = svc.PutItem(input)
					if err != nil {
						log.Fatalf("Error calling PutItem: %s", err)
					}

					fmt.Println("Data added to table " + tableName)*/
			}
		}
		time.Sleep(3600 * time.Second)
	}
}
