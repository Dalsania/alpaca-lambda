package main

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/shopspring/decimal"
)

type MClient struct {
}
type AClient struct {
}

var forty = decimal.NewFromFloat(40)
var twenty = decimal.NewFromFloat(20)
var four = decimal.NewFromFloat(4)

// var cat = decimal.NewFromFloat(40)
// var cc = &cat

var GetClockMock func() (*alpaca.Clock, error)
var PlaceOrderMock func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
var GetSnapshotMock func(symbol string) (*marketdata.Snapshot, error)
var MClientMock func() Marketer
var AClientMock func() Alpacer

func initMocks() {
	GetClockMock = func() (*alpaca.Clock, error) {
		return &alpaca.Clock{IsOpen: true}, nil
	}
	PlaceOrderMock = func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
		return &alpaca.Order{}, nil
	}
	GetSnapshotMock = func(symbol string) (*marketdata.Snapshot, error) {
		return &marketdata.Snapshot{LatestTrade: &marketdata.Trade{Price: 10}, PrevDailyBar: &marketdata.Bar{Close: 15}}, nil
	}
	MClientMock = func() Marketer {
		return MClient{}
	}
	AClientMock = func() Alpacer {
		return AClient{}
	}
	initAlpacaClient = AClientMock
	initMarketClient = MClientMock
}

func (m MClient) GetSnapshot(symbol string) (*marketdata.Snapshot, error) {
	return GetSnapshotMock(symbol)
}

func (a AClient) GetClock() (*alpaca.Clock, error) {
	return GetClockMock()
}
func (a AClient) PlaceOrder(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
	return PlaceOrderMock(req)
}

func TestAlpaca(t *testing.T) {
	cases := []struct {
		Name        string
		MClient     func() Marketer
		AClient     func() Alpacer
		GetClock    func() (*alpaca.Clock, error)
		GetSnapshot func(symbol string) (*marketdata.Snapshot, error)
		PlaceOrder  func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
		ExpErr      bool
	}{
		{
			Name:   "Happy Case",
			ExpErr: false,
		},
		{
			Name: "Market Call Fails",
			GetClock: func() (*alpaca.Clock, error) {
				return &alpaca.Clock{}, errors.New("Failed to get clock")
			},
			ExpErr: true,
		},
		{
			Name: "Market is Closed",
			GetClock: func() (*alpaca.Clock, error) {
				return &alpaca.Clock{IsOpen: false}, nil
			},
			ExpErr: true,
		},
		{
			Name: "Get Snapshot fails",
			GetSnapshot: func(symbol string) (*marketdata.Snapshot, error) {
				return &marketdata.Snapshot{}, errors.New("Failed to get snapshot")
			},
			ExpErr: true,
		},
	}

	for _, tc := range cases {
		initMocks()
		if tc.GetClock != nil {
			GetClockMock = tc.GetClock
		}
		if tc.GetSnapshot != nil {
			GetSnapshotMock = tc.GetSnapshot
		}

		_, err := HandleRequest()
		if err == nil && tc.ExpErr {
			t.Errorf("Expected error but got %s", err)
		}
		if err != nil && !tc.ExpErr {
			t.Errorf("Did not expect error but got %s", err)
		}
	}

}
func TestSnapshot(t *testing.T) {

	cases := []struct {
		Name        string
		GetSnapshot func(symbol string) (*marketdata.Snapshot, error)
		PlaceOrder  func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error)
		Res         *alpaca.Order
	}{
		{
			Name: "Buy 40",
			GetSnapshot: func(symbol string) (*marketdata.Snapshot, error) {
				return &marketdata.Snapshot{LatestTrade: &marketdata.Trade{Price: 388.285}, PrevDailyBar: &marketdata.Bar{Close: 3389.00}}, nil
			},
			PlaceOrder: func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
				n := req.Notional
				s := req.AssetKey
				return &alpaca.Order{Symbol: *s, Notional: n}, nil
			},
			Res: &alpaca.Order{Symbol: "VOO", Notional: &forty},
		},
		{
			Name: "Buy Nothing",
			GetSnapshot: func(symbol string) (*marketdata.Snapshot, error) {
				return &marketdata.Snapshot{LatestTrade: &marketdata.Trade{Price: 388.285}, PrevDailyBar: &marketdata.Bar{Close: 389.00}}, nil
			},
			Res: &alpaca.Order{},
		},
		{
			Name: "Buy 20",
			GetSnapshot: func(symbol string) (*marketdata.Snapshot, error) {
				return &marketdata.Snapshot{LatestTrade: &marketdata.Trade{Price: 388.285}, PrevDailyBar: &marketdata.Bar{Close: 400.00}}, nil
			},
			PlaceOrder: func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
				n := req.Notional
				s := req.AssetKey
				return &alpaca.Order{Symbol: *s, Notional: n}, nil
			},
			Res: &alpaca.Order{Symbol: "VOO", Notional: &twenty},
		},
		{
			Name: "Buy 4",
			GetSnapshot: func(symbol string) (*marketdata.Snapshot, error) {
				return &marketdata.Snapshot{LatestTrade: &marketdata.Trade{Price: 388.285}, PrevDailyBar: &marketdata.Bar{Close: 396.00}}, nil
			},
			PlaceOrder: func(req alpaca.PlaceOrderRequest) (*alpaca.Order, error) {
				n := req.Notional
				s := req.AssetKey
				return &alpaca.Order{Symbol: *s, Notional: n}, nil
			},
			Res: &alpaca.Order{Symbol: "VOO", Notional: &four},
		},
	}
	for _, tc := range cases {
		initMocks()
		if tc.GetSnapshot != nil {
			GetSnapshotMock = tc.GetSnapshot
		}
		if tc.PlaceOrder != nil {
			PlaceOrderMock = tc.PlaceOrder
		}
		res, _ := HandleRequest()

		dq := cmp.Equal(*tc.Res, *res)

		if dq == false && (&alpaca.Order{}) == res {
			t.Log(tc.Name)
			t.Errorf("Shouldn't have bought anything but you did")
		}
		if (alpaca.Order{}) != *res {

			c := *res.Notional
			if c.String() != tc.Res.Notional.String() {
				t.Errorf("Order was not placed correctly: Bought $ %s worth , Expected to buy $ %s worth of stock", res.Notional, tc.Res.Notional)
			}
		}
	}
}
