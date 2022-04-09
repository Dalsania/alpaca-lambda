package main

import(
	//"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

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

func main(){
	// err := godotenv.Load("../.env")
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	//   }
	   APCA_API_KEY_ID := os.Getenv("APCA_API_KEY_ID")
	   APCA_API_SECRET_KEY := os.Getenv("APCA_API_SECRET_KEY")

	   fmt.Println(APCA_API_KEY_ID)
	   fmt.Println(APCA_API_SECRET_KEY)
	  marketClient := marketdata.NewClient(marketdata.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,
		BaseURL: "https://data.alpaca.markets",
	})
	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		ApiKey:    APCA_API_KEY_ID,
		ApiSecret: APCA_API_SECRET_KEY,,
		BaseURL: "https://paper-api.alpaca.markets",
		})
	 
	
	 var not decimal.Decimal
	 snapshot,err:= marketClient.GetSnapshot("TWTR")

	 if err != nil {
		 log.Printf("Failed to retrieve Snapshot Data: %p",err)
	 }


	 var prevDailyClose = snapshot.PrevDailyBar.Close
	 latestTrade,err := marketClient.GetLatestTrade("TWTR")
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
		not = decimal.NewFromInt(4)
	}
	if perChange < -0.015 && perChange > -0.02 {
		fmt.Println("-0.015 to -0.020: buy 8 dollars")
		not = decimal.NewFromInt(8)
	}
	if perChange < -0.02 && perChange > -0.04 {
		fmt.Println("-0.02 to -0.04: buy 10 dollars")
		not = decimal.NewFromInt(10)
	}
	if perChange < -0.05  {
		fmt.Println("losses greater than -0.05: buy 12 dollars")
		not = decimal.NewFromInt(12)
	}
	orderStruct := createOrder(not, "TWTR")
	_,err = alpacaClient.PlaceOrder(orderStruct)
	if err != nil{
		log.Printf("Failed to place order: %s", err)
	}
}