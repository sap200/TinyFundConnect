package types

import "encoding/json"

type DepositRequest struct {
	PoolId              string  `json:"poolId"`
	SuccessRedirectLink string  `json:"successRedirectLink"`
	ErrorRedirectLink   string  `json:"errorRedirectLink"`
	Amount              float64 `json:"amount"`
	UserEmail           string  `json:"userEmail"`
}

type DepositResponse struct {
	PoolId       string  `json:"poolId"`
	UserEmail    string  `json:"userEmail"`
	Amount       float64 `json:"amount"`
	CheckoutLink string  `json:"checkoutLink"`
	CheckoutId   string  `json:"checkoutId"`
	MerchantId   string  `json:"merchantId"`
}

func (m DepositResponse) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

const DEPOSIT = "Deposit"
const WITHDRAWL = "Withdrawl"

type Transaction struct {
	PoolId     string  `json:"poolId"`
	UserEmail  string  `json:"userEmail"`
	CheckoutId string  `json:"checkoutId"`
	Amount     float64 `json:"amount"`
	Paid       bool    `json:"paid"`
	TxnType    string  `json:"txnType"`
	PayoutId   string  `json:"payoutId"`
}

func (m Transaction) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}

type WithdrawlRequest struct {
	FirstName          string  `json:"firstName"`
	LastName           string  `json:"lastName"`
	BankAccountNumber1 string  `json:"bankAccountNumber1"`
	BankAccountNumber2 string  `json:"bankAccountNumber2"`
	BankBranchCode     string  `json:"bankBranchCode"`
	PhoneNumber        string  `json:"phoneNumber"`
	Email              string  `json:"email"`
	Amount             float64 `json:"amount"`
}

type WithdrawlResponse struct {
	Status   string  `json:"status"`
	PayoutId string  `json:"payoutId"`
	Amount   float64 `json:"amount"`
}

func (m WithdrawlResponse) EncodeToMap() (map[string]interface{}, error) {
	// Create a custom decoder configuration
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(inrec, &inInterface)
	return inInterface, nil
}
