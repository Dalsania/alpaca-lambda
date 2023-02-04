package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/shopspring/decimal"
)

//API KEYS
var APCA_API_KEY_ID = os.Getenv("APCA_API_KEY_ID")
var APCA_API_SECRET_KEY = os.Getenv("APCA_API_SECRET_KEY")

func createOrder(not decimal.Decimal, symbol string) alpaca.PlaceOrderRequest {
	var marketOrder = alpaca.PlaceOrderRequest{
		AssetKey:    &symbol,
		Notional:    &not,
		Side:        "buy",
		Type:        "market",
		TimeInForce: "day",
	}
	return marketOrder
}

//marketdata.Client

type Marketer interface {
	GetSnapshot(symbol string) (*marketdata.Snapshot, error)
}
type Alpacer interface {
	GetClock() (*alpaca.Clock, error)
	PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
}

var initMarketClient = func() Marketer {
	var marketClient Marketer
	marketClient = marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,
		BaseURL:   "https://data.alpaca.markets",
	})
	return marketClient

}
var initAlpacaClient = func() Alpacer {
	var alpacaClient Alpacer
	alpacaClient = alpaca.NewClient(alpaca.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,
		BaseURL:   "https://api.alpaca.markets",
		//BaseURL:   "https://paper-api.alpaca.markets",
	})
	return alpacaClient
}

//ctx context.Context
func HandleRequest() (*alpaca.Order, error) {

	mc := initMarketClient()
	ac := initAlpacaClient()
	clock, err := ac.GetClock()

	if err != nil {
		msg := "Failed to get current market timestamp"
		log.Println(msg)
		return &alpaca.Order{}, errors.New(msg)
	}
	s := strconv.FormatBool(clock.IsOpen)
	log.Printf("Market is open: %s \n", s)
	if clock.IsOpen == false {
		msg := "Market is not open"
		log.Println(msg)
		return &alpaca.Order{}, nil
	}
	var not decimal.Decimal
	snapshot, err := mc.GetSnapshot("VOO")

	if err != nil {
		msg := "Failed to retrieve Snapshot Data:"
		log.Println(msg)
		return &alpaca.Order{}, errors.New(msg)
	}

	var prevDailyClose = snapshot.PrevDailyBar.Close
	//marketClient.GetLatestTrade("VOO")
	//latestTrade := snapshot.LatestTrade.Price
	// if err != nil {
	// 	log.Printf("Failed to get latetest trade for stock: %p", err)
	// }
	var lastTradePrice = snapshot.LatestTrade.Price

	fmt.Printf("Previous Daily Close: %f \n", prevDailyClose)
	fmt.Printf("Last Trading Price of stock %f \n", lastTradePrice)

	var diff = lastTradePrice - prevDailyClose
	fmt.Printf("The diff between the latest trade of the stock and yesterdays close: %f \n", diff)
	var perChange = diff / prevDailyClose
	fmt.Printf("The percent change: %f \n", perChange)

	if perChange > -0.015 {
		log.Println("? to -0.015: Not Buying anything")
		return &alpaca.Order{}, nil
	}
	if perChange < -0.015 && perChange > -0.02 {
		log.Println("-0.015 to -0.020: buy 4 dollars")
		not = decimal.NewFromInt(4)
	}
	if perChange < -0.02 && perChange > -0.04 {
		log.Println("-0.02 to -0.04: buy 20 dollars")
		not = decimal.NewFromInt(20)
	}
	if perChange < -0.05 {
		log.Println("losses greater than -0.05: buy 40 dollars")
		not = decimal.NewFromInt(40)
	}

	orderStruct := createOrder(not, "VOO")
	res, err := ac.PlaceOrder(orderStruct)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	}
	return res, nil
}

func main() {
	lambda.Start(HandleRequest)

}
