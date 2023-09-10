package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sap200/TinyFundConnect/secret"
)

func GetUnsignedBinanceRequest(requestType, path string, paramsMap map[string]string) (*http.Request, error) {
	params := url.Values{}

	for k, v := range paramsMap {
		params.Set(k, v)
	}

	req, err := http.NewRequest(requestType, path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.URL.RawQuery = params.Encode()

	return req, nil

}

func GetSignedBinanceRequest(requestType, path string, paramsMap map[string]string) (*http.Request, error) {
	params := url.Values{}

	for k, v := range paramsMap {
		params.Set(k, v)
	}

	// get and set the timestamp
	timestamp := fmt.Sprint(time.Now().UnixNano() / int64(time.Millisecond))
	params.Set("timestamp", timestamp)

	// Create the query string from the parameters
	queryString := params.Encode()
	fmt.Println(queryString)

	// Create the payload by combining the HTTP method, endpoint, and query string
	//payload := fmt.Sprintf("%s\n%s\n%s", requestType, path, queryString)
	// Generate the HMAC-SHA256 signature
	signature := generateSignature(secret.BINANCE_SPOT_TEST_NET_SECRET_KEY, queryString)
	params.Set("signature", signature)

	headers := map[string]string{
		"X-MBX-APIKEY": secret.BINANCE_SPOT_TEST_NET_API_KEY,
	}

	// Create an HTTP request
	req, err := http.NewRequest(requestType, path, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Set the query parameters
	req.URL.RawQuery = params.Encode()

	return req, nil

}

func generateSignature(secretKey, payload string) string {
	key := []byte(secretKey)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(payload))
	return hex.EncodeToString(h.Sum(nil))
}
