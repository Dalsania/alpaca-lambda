package main

import (
	"log"
	"github.com/shopspring/decimal"
	"os"
	"fmt"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/aws/aws-lambda-go/lambda"
	

	"context"
)
//API KEYS
var APCA_API_KEY_ID = os.Getenv("APCA_API_KEY_ID")
var APCA_API_SECRET_KEY = os.Getenv("APCA_API_SECRET_KEY")


func createOrder(not decimal.Decimal, symbol string) alpaca.PlaceOrderRequest{
	  var marketOrder = alpaca.PlaceOrderRequest{
	   AssetKey: &symbol,
	   Notional:&not,
	   Side:"buy",
	   Type:"market",
	   TimeInForce:"day",
	   }
	   return marketOrder
}
func HandleRequest(ctx context.Context) {

	marketClient := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,
		BaseURL: "https://data.alpaca.markets",
	})
	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,
		BaseURL: "https://paper-api.alpaca.markets",
		})

	 var not decimal.Decimal
	 snapshot,err:= marketClient.GetSnapshot("VOO")

	 if err != nil {
		 log.Printf("Failed to retrieve Snapshot Data: %p",snapshot)
	 }


	 var prevDailyClose = snapshot.PrevDailyBar.Close
	 latestTrade,err := marketClient.GetLatestTrade("VOO")
	 if err != nil{
		 log.Printf("Failed to get latetest trade for stock: %p", err)
	 }
	 var lastTradePrice = latestTrade.Price

	 fmt.Printf("Previous Daily Close: %f \n", prevDailyClose)
	 fmt.Printf("Last Trading Price of stock %f \n",lastTradePrice)
	 
	 var diff = lastTradePrice - prevDailyClose
	 fmt.Printf("The diff between the latest trade of the stock and yesterdays close: %f \n",diff)
	 var perChange = diff/prevDailyClose
	 fmt.Printf("The percent change: %f \n",perChange)
	
	if perChange < 0.02 && perChange > -0.015{
		fmt.Println("-0.015 to 0.020: buy 4 dollars")
		not = decimal.NewFromInt(1)
	}
	if perChange < -0.015 && perChange > -0.02 {
		fmt.Println("-0.015 to -0.020: buy 8 dollars")
		not = decimal.NewFromInt(4)
	}
	if perChange < -0.02 && perChange > -0.04 {
		fmt.Println("-0.02 to -0.04: buy 10 dollars")
		not = decimal.NewFromInt(8)
	}
	if perChange < -0.05  {
		fmt.Println("losses greater than -0.05: buy 12 dollars")
		not = decimal.NewFromInt(12)
	}
	orderStruct := createOrder(not, "VOO")
	_,err = alpacaClient.PlaceOrder(orderStruct)
	if err != nil{
		log.Printf("Failed to place order: %s", err)
	}

	}

func main() {
	lambda.Start(HandleRequest)
	// fmt.Printf("%+v\n", *acct)
}