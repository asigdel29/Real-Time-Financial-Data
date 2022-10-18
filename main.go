package main

import (
	"encoding/json"
	"fmt"
	"github.com/JamesPEarly/loggly"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Item struct {
	StockSymbol          string
	Time                 string
	Price                float64 `json:"Price"`
	DollarChange         float64 `json:"DollarChange"`
	PercentChange        float64 `json:"PercentChange"`
	PreviousClose        float64 `json:"PreviousClose"`
	Open                 float64 `json:"Open"`
	BidPrice             float64 `json:"BidPrice"`
	BidQuantity          int     `json:"BidQuantity"`
	AskPrice             float64 `json:"AskPrice"`
	AskQuantity          int     `json:"AskQuantity"`
	DayRangeLow          float64 `json:"DayRangeLow"`
	DayRangeHigh         float64 `json:"DayRangeHigh"`
	YearRangeLow         float64 `json:"YearRangeLow"`
	YearRangeHigh        float64 `json:"YearRangeHigh"`
	Volume               int     `json:"Volume"`
	AverageVolume        int     `json:"AverageVolume"`
	MarketCap            float64 `json:"MarketCap"`
	Beta                 float64 `json:"Beta"`
	PriceEarningsRatio   float64 `json:"PriceEarningsRatio"`
	EarningsPerShare     float64 `json:"EarningsPerShare"`
	EarningsDate         string  `json:"EarningsDate"`
	ForwardDividend      float64 `json:"ForwardDividend"`
	ForwardDividendYield float64 `json:"ForwardDividendYield"`
	ExDividendDate       string  `json:"ExDividendDate"`
	YearTargetEstimate   float64 `json:"YearTargetEstimate"`
	QueriedSymbol        string  `json:"QueriedSymbol"`
}

type Response struct {
	Summary Summary `json:"Summary"`
}

type Summary struct {
	Name                 string    `json:"Name"`
	StockSymbol          string    `json:"StockSymbol"`
	Price                float64   `json:"Price"`
	DollarChange         float64   `json:"DollarChange"`
	PercentChange        float64   `json:"PercentChange"`
	PreviousClose        float64   `json:"PreviousClose"`
	Open                 float64   `json:"Open"`
	BidPrice             float64   `json:"BidPrice"`
	BidQuantity          int       `json:"BidQuantity"`
	AskPrice             float64   `json:"AskPrice"`
	AskQuantity          int       `json:"AskQuantity"`
	DayRangeLow          float64   `json:"DayRangeLow"`
	DayRangeHigh         float64   `json:"DayRangeHigh"`
	YearRangeLow         float64   `json:"YearRangeLow"`
	YearRangeHigh        float64   `json:"YearRangeHigh"`
	Volume               int       `json:"Volume"`
	AverageVolume        int       `json:"AverageVolume"`
	MarketCap            float64   `json:"MarketCap"`
	Beta                 float64   `json:"Beta"`
	PriceEarningsRatio   float64   `json:"PriceEarningsRatio"`
	EarningsPerShare     float64   `json:"EarningsPerShare"`
	EarningsDate         string    `json:"EarningsDate"`
	ForwardDividend      float64   `json:"ForwardDividend"`
	ForwardDividendYield float64   `json:"ForwardDividendYield"`
	ExDividendDate       string    `json:"ExDividendDate"`
	YearTargetEstimate   float64   `json:"YearTargetEstimate"`
	QueriedSymbol        string    `json:"QueriedSymbol"`
	DataCollectedOn      time.Time `json:"DataCollectedOn"`
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
		for i := 1; i <= 10; i++ {
			for _, element := range stocks {

				weekday := time.Now().Weekday()
				daycheck := (int(weekday))
				if daycheck == 7 || daycheck == 0 {
					fmt.Println("Market Closed")
					os.Exit(3)
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

					var response Response
					json.Unmarshal(responseBytes, &response)
					err = client.EchoSend("info", "Request succeeded")
					if err != nil {
						fmt.Println("err:", err)
					}
					formattedData, _ := json.MarshalIndent(response, "", "    ")
					fmt.Println(string(formattedData))

					var respSize string = strconv.Itoa(len(responseBytes))
					logErr := client.EchoSend("info", "Successful data collection of size: "+respSize)
					if logErr != nil {
						fmt.Println("err: ", logErr)
					}

					sess, err := session.NewSession(&aws.Config{
						Region: aws.String("us-east-1")},
					)

					if err != nil {
						log.Fatalf("Error initializing AWS: %s", err)
					}

					svc := dynamodb.New(sess)

					var item Item
					item.Price = response.Summary.Price
					item.DollarChange = response.Summary.DollarChange
					item.PercentChange = response.Summary.PercentChange
					item.PreviousClose = response.Summary.PreviousClose
					item.Open = response.Summary.Open
					item.BidPrice = response.Summary.Open
					item.BidQuantity = response.Summary.BidQuantity
					item.AskPrice = response.Summary.AskPrice
					item.AskQuantity = response.Summary.AskQuantity
					item.DayRangeLow = response.Summary.DayRangeLow
					item.DayRangeHigh = response.Summary.DayRangeHigh
					item.YearRangeLow = response.Summary.YearRangeLow
					item.YearRangeHigh = response.Summary.YearRangeHigh
					item.Volume = response.Summary.Volume
					item.AverageVolume = response.Summary.AverageVolume
					item.MarketCap = response.Summary.MarketCap
					item.Beta = response.Summary.Beta
					item.PriceEarningsRatio = response.Summary.PriceEarningsRatio
					item.EarningsPerShare = response.Summary.EarningsPerShare
					item.EarningsDate = response.Summary.EarningsDate
					item.ForwardDividendYield = response.Summary.ForwardDividendYield
					item.ForwardDividend = response.Summary.ForwardDividend
					item.ExDividendDate = response.Summary.ExDividendDate
					item.YearTargetEstimate = response.Summary.YearTargetEstimate
					item.QueriedSymbol = response.Summary.QueriedSymbol
					item.Time = time.Now().Format(time.RFC3339)
					item.StockSymbol = response.Summary.StockSymbol

					av, err := dynamodbattribute.MarshalMap(item)
					if err != nil {
						log.Fatalf("Error marshalling %s", err)
					}

					tableName := "asigdel-topstocks"
					input := &dynamodb.PutItemInput{
						Item:      av,
						TableName: aws.String(tableName),
					}

					_, err = svc.PutItem(input)
					if err != nil {
						log.Fatalf("Error calling PutItem: %s", err)
					}

					fmt.Println("Data added to table " + tableName)
				}
			}
		}
		time.Sleep(3600 * time.Second)
	}
}
