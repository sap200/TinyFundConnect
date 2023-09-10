package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/sap200/TinyFundConnect/db"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"github.com/sap200/TinyFundConnect/utils"
)

func GetAllOrdersByPoolId(paramsMap map[string]string, poolId string) ([]map[string]interface{}, error) {
	path := "/api/v3/allOrders"

	req, err := GetSignedBinanceRequest(http.MethodGet, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

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

	// get repo of pool
	repo := db.New()
	exists, err := repo.DoesRecordExists(secret.ORDERS_COLLECTION, poolId)
	if err != nil {
		return nil, err
	}

	var oe *map[string]types.OrderEarns
	myOwnMap := make(map[string]types.OrderEarns)
	oe = &myOwnMap

	if exists {
		oe, err = repo.GetOrderEarnsByPoolId(secret.ORDERS_COLLECTION, poolId)
		if err != nil {
			return nil, err
		}
	}

	//fmt.Println(*oe)

	//fmt.Println(string(body))

	var data []map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	nm := []map[string]interface{}{}

	for i := 0; i < len(data); i++ {
		cid, _ := data[i]["clientOrderId"].(string)
		if cid == poolId {
			// check in map if orderid exists
			jsonBytes, _ := json.Marshal(data[i]["orderId"])
			oid := string(jsonBytes)
			//fmt.Println(oid)
			val, ok := (*oe)[oid]
			if !ok {
				nm = append(nm, data[i])
			} else {
				mymap := data[i]
				mymap["totalEarned"] = val.TotalEarning / 83.12
				nm = append(nm, mymap)
			}
		}
	}

	return nm, nil

}

func GetOpenOrdersByPoolId(paramsMap map[string]string, poolId string) ([]map[string]interface{}, error) {
	path := "/api/v3/openOrders"

	req, err := GetSignedBinanceRequest(http.MethodGet, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

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

	nm := []map[string]interface{}{}

	for i := 0; i < len(data); i++ {
		cid, _ := data[i]["clientOrderId"].(string)
		if cid == poolId {
			nm = append(nm, data[i])
		}
	}

	return nm, nil
}

func CancelOrderByOrderId(paramsMap map[string]string, earned float64, emailId, poolId string) (map[string]interface{}, error) {

	path := "/api/v3/order"

	req, err := GetSignedBinanceRequest(http.MethodDelete, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	data["totalProfitEarned"] = earned

	// run a go routine to save orders to firebase db
	floatVal := earned * 83.12
	orderId := paramsMap["orderId"]
	keyPoolId := poolId
	go func() {
		SaveEarnedFromOrder(floatVal, orderId, keyPoolId)
	}()

	// run a go routine to apportion the amount
	go func() {
		ApportionProfitBetweenPoolMembers(floatVal, keyPoolId)
	}()

	// TODO:: Make a log that following order is cancelled by following email user in following pool

	return data, nil
}

func PlaceLimitOrder(paramsMap map[string]string, email, poolId string) (map[string]interface{}, error) {

	path := "/api/v3/order"

	req, err := GetSignedBinanceRequest(http.MethodPost, secret.BINANCE_SPOT_TEST_NET_ENDPOINT+path, paramsMap)
	if err != nil {
		return nil, err
	}

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

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// TODO:: Make a log that following Limit order is Created by following email user in following pool

	return data, nil
}

func ApportionProfitBetweenPoolMembers(amount float64, poolId string) error {

	fmt.Println("Inside apportionment")

	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, poolId)
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("Pool record not found")
	}

	fmt.Println("Getting user in apportionment")

	pd, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		return err
	}

	var pool types.Pool
	json.Unmarshal(utils.GetBytesFromInterface(pd), &pool)

	totalBalance := pool.PoolBalance

	if pool.MembersList == nil {
		return errors.New("Unexpected pool member list is empty")
	}

	if totalBalance == 0 {
		return errors.New("Pool Balance is 0, This function can be executed only when total pool balance is > 0")
	}

	for i := 0; i < len(pool.MembersList); i++ {
		apportionedAmount := (pool.MembersList[i].Balance / totalBalance) * amount
		pool.MembersList[i].Balance += apportionedAmount
		fmt.Println(apportionedAmount)
		// get db record and update user record here
		exists, err := repo.DoesRecordExists(secret.USER_COLLECTIONS, pool.MembersList[i].Email)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}

		ud, err := repo.GetRecordDetails(secret.USER_COLLECTIONS, pool.MembersList[i].Email)
		if err != nil {
			return err
		}
		var userDetails types.User
		json.Unmarshal(utils.GetBytesFromInterface(ud), &userDetails)
		userDetails.Balance += apportionedAmount
		um, _ := userDetails.EncodeToMap()
		repo.Save(secret.USER_COLLECTIONS, pool.MembersList[i].Email, um)
	}

	// recalculate total balance
	tb := float64(0)
	for i := 0; i < len(pool.MembersList); i++ {
		tb += pool.MembersList[i].Balance
	}

	pool.PoolBalance = tb

	pm, _ := pool.EncodeToMap()

	fmt.Println("Amount Apportioned")
	fmt.Println(pool)

	repo.Save(secret.POOL_COLLECTION, poolId, pm)

	return nil
}

func SaveEarnedFromOrder(amount float64, orderId string, poolId string) error {
	repo := db.New()

	e, err := repo.DoesRecordExists(secret.ORDERS_COLLECTION, poolId)
	if err != nil {
		return err
	}
	if !e {

		oe := map[string]types.OrderEarns{}
		oe[orderId] = types.OrderEarns{
			OrderId:      orderId,
			PoolId:       poolId,
			TotalEarning: amount,
		}

		d, _ := json.Marshal(oe)
		var m map[string]interface{}
		json.Unmarshal(d, &m)
		repo.Save(secret.ORDERS_COLLECTION, poolId, m)
		return nil
	}

	oearns, err := repo.GetOrderEarnsByPoolId(secret.ORDERS_COLLECTION, poolId)
	(*oearns)[orderId] = types.OrderEarns{
		OrderId:      orderId,
		PoolId:       poolId,
		TotalEarning: amount,
	}

	d, _ := json.Marshal(oearns)
	var m map[string]interface{}
	json.Unmarshal(d, &m)
	repo.Save(secret.ORDERS_COLLECTION, poolId, m)

	return nil
}

func GetTotalEarnedFromOrderId(orderId, poolId string) (float64, error) {
	repo := db.New()

	e, err := repo.DoesRecordExists(secret.ORDERS_COLLECTION, poolId)
	if err != nil {
		return 0, err
	}

	if !e {
		return 0, nil
	}

	m, err := repo.GetOrderEarnsByPoolId(secret.ORDERS_COLLECTION, poolId)
	if err != nil {
		return 0, err
	}

	a, ok := (*m)[orderId]
	if !ok {
		return 0, nil
	}

	return a.TotalEarning, nil
}
