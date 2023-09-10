package binance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
)

func GetMarketData(paramsMap map[string]string) (*[]types.KLineData, error) {
	path := "/api/v3/uiKlines"

	req, err := GetUnsignedBinanceRequest(http.MethodGet, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	var listOfCandles []types.KLineData
	listOfCandles = []types.KLineData{}

	for i := 0; i < len(data); i++ {
		a := data[i]
		openTime := a[0]
		closeTime := a[6]
		numberOfTrades := a[8]
		openPrice, _ := a[1].(string)
		highPrice, _ := a[2].(string)
		lowPrice, _ := a[3].(string)
		closePrice, _ := a[4].(string)
		volume, _ := a[5].(string)
		quoteAssetVolume, _ := a[7].(string)
		takerBuyBaseAssetVolume, _ := a[9].(string)
		TakerBuyQuoteAssetVolume, _ := a[10].(string)
		unusedField, _ := a[11].(string)

		c := types.KLineData{
			OpenTime:                 openTime,
			OpenPrice:                openPrice,
			HighPrice:                highPrice,
			LowPrice:                 lowPrice,
			ClosePrice:               closePrice,
			Volume:                   volume,
			CloseTime:                closeTime,
			QuoteAssetVolume:         quoteAssetVolume,
			NumberOfTrades:           numberOfTrades,
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolume,
			TakerBuyQuoteAssetVolume: TakerBuyQuoteAssetVolume,
			UnusedField:              unusedField,
		}

		listOfCandles = append(listOfCandles, c)
	}

	return &listOfCandles, nil
}

func GetOrderBook(paramsMap map[string]string) (*types.OrderBook, error) {
	path := "/api/v3/depth"

	req, err := GetUnsignedBinanceRequest(http.MethodGet, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	//fmt.Println(data)

	bidList, _ := data["bids"].([]interface{})
	askList, _ := data["asks"].([]interface{})

	bl := []types.Bid{}
	al := []types.Ask{}

	for i := 0; i < len(bidList); i++ {
		price := bidList[i].([]interface{})[0].(string)
		qty := bidList[i].([]interface{})[1].(string)
		b := types.Bid{
			Qty:   qty,
			Price: price,
		}
		bl = append(bl, b)

		price1 := askList[i].([]interface{})[0].(string)
		qty1 := askList[i].([]interface{})[1].(string)
		a := types.Ask{
			Qty:   qty1,
			Price: price1,
		}
		al = append(al, a)
	}

	ob := types.OrderBook{
		Asks:       al,
		Bids:       bl,
		LastUpdate: data["lastUpdateId"],
	}

	return &ob, nil

}

func GetRecentTradingList(paramsMap map[string]string) ([]map[string]interface{}, error) {
	path := "/api/v3/trades"

	req, err := GetUnsignedBinanceRequest(http.MethodGet, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

	// Create an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
