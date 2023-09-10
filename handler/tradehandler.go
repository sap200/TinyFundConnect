package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/binance"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangealogger"
	"github.com/sap200/TinyFundConnect/types"
)

func GetCandleDatHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	interval := c.Query("interval")
	limit := c.Query("limit")

	paramsMap := map[string]string{
		"symbol":   symbol,
		"interval": interval,
		"limit":    limit,
	}

	kd, err := binance.GetMarketData(paramsMap)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, *kd)
}

func RetrieveOrderBookHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	limit := c.Query("limit")

	params := map[string]string{
		"symbol": symbol,
		"limit":  limit,
	}

	ob, err := binance.GetOrderBook(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, ob)
}

func GetRecentTradeHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	limit := c.Query("limit")

	params := map[string]string{
		"symbol": symbol,
		"limit":  limit,
	}

	m, err := binance.GetRecentTradingList(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, m)
}

func GetTradeOrderByPoolIdHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	poolId := c.Query("poolId")

	params := map[string]string{
		"symbol": symbol,
	}

	ao, err := binance.GetAllOrdersByPoolId(params, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, ao)

}

func GetTradeOpenOrdersByPoolId(c *gin.Context) {
	symbol := c.Query("symbol")
	poolId := c.Query("poolId")

	params := map[string]string{
		"symbol": symbol,
	}

	ao, err := binance.GetOpenOrdersByPoolId(params, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, ao)
}

func CancelOpenOrderTradeHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	orderId := c.Query("orderId")
	poolId := c.Query("poolId")
	email := c.Query("email")
	earned := c.Query("earned")
	earnedAmount, _ := strconv.ParseFloat(earned, 64)

	params := map[string]string{
		"symbol":  symbol,
		"orderId": orderId,
	}

	co, err := binance.CancelOrderByOrderId(params, earnedAmount, email, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	event := pangealogger.New(
		"User with email "+email+" closed an order in pool with id "+poolId,
		"ORDER_CLOSE_EVENT",
		"",
		email,
		poolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, co)
}

func CreateLimitOrderTradeHandler(c *gin.Context) {
	symbol := c.Query("symbol")
	poolId := c.Query("poolId")
	email := c.Query("email")
	side := c.Query("side")
	quantity := c.Query("quantity")
	price := c.Query("price")

	paramsMap := map[string]string{
		"symbol":           symbol,
		"newClientOrderId": poolId,
		"side":             side,
		"type":             "LIMIT",
		"timeInForce":      "GTC",
		"quantity":         quantity,
		"price":            price,
	}

	cro, err := binance.PlaceLimitOrder(paramsMap, email, poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "FAILURE",
			"code":    e.TRADE_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	event := pangealogger.New(
		"User with email "+email+" created a LIMIT order in pool with id "+poolId,
		"ORDER_CREATION_EVENT",
		"",
		email,
		poolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, cro)
}

func GetAllSymbolTradeHandler(c *gin.Context) {
	a := types.ListAllSymbols()

	c.JSON(http.StatusOK, a)
}
