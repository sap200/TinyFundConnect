package types

import "encoding/json"

type KLineData struct {
	OpenTime                 interface{} `json:"openTime"`
	OpenPrice                string      `json:"openPrice"`
	HighPrice                string      `json:"highPrice"`
	LowPrice                 string      `json:"lowPrice"`
	ClosePrice               string      `json:"closePrice"`
	Volume                   string      `json:"volume"`
	CloseTime                interface{} `json:"closeTime"`
	QuoteAssetVolume         string      `json:"quoteAssetVolume"`
	NumberOfTrades           interface{} `json:"numberOfTrades"`
	TakerBuyBaseAssetVolume  string      `json:"takerBuyBaseAssetVolume"`
	TakerBuyQuoteAssetVolume string      `json:"takerBuyQuoteAssetVolume"`
	UnusedField              string      `json:"unusedField"`
}

type Ask struct {
	Qty   string `json:"qty`
	Price string `json:"price"`
}

type Bid struct {
	Qty   string `json:"qty"`
	Price string `json:"price"`
}

type OrderBook struct {
	LastUpdate interface{} `json:"lastUpdated"`
	Bids       []Bid       `json:"bids"`
	Asks       []Ask       `json:"asks"`
}

func ListAllSymbols() []string {
	return []string{
		"BNBUSDT",
		"BTCUSDT",
		"ETHBTC",
		"ETHUSDT",
		"BNBBTC",
	}
}

type OrderEarns struct {
	PoolId       string  `json:"poolId"`
	TotalEarning float64 `json:"totalEarning"`
	OrderId      string  `json:"orderId"`
}

func (v OrderEarns) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
