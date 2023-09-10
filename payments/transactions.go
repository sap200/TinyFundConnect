package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sap200/TinyFundConnect/db"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"github.com/sap200/TinyFundConnect/utils"
)

func StartDeposit(poolId, email, successRedirect, errorRedirect string, amount float64) (*types.DepositResponse, error) {
	if amount <= 0 {
		return nil, errors.New("Amount has to be >= 0")
	}

	// create a rapyd transaction
	dr, err := sendRapydPaymentRequest(poolId, email, successRedirect, errorRedirect, amount)
	if err != nil {
		return nil, err
	}

	fmt.Println(*dr)

	return dr, nil
}

func sendRapydPaymentRequest(poolId, email, successRedirect, errorRedirect string, amount float64) (*types.DepositResponse, error) {

	type Md struct {
		MerchantDefined bool `json:"merchant_defined"`
	}

	merchantReferenceId := email + "_" + poolId + "_" + uuid.New().String()

	reqBody := struct {
		Amount                      float64 `json:"amount"`
		CardholderPreferredCurrency bool    `json:"cardholder_preferred_currency"`
		CompletePaymentURL          string  `json:"complete_payment_url"`
		Country                     string  `json:"country"`
		Currency                    string  `json:"currency"`
		ErrorPaymentURL             string  `json:"error_payment_url"`
		Language                    string  `json:"language"`
		MerchantReferenceID         string  `json:"merchant_reference_id"`
		Metadata                    Md      `json:"metadata"`
	}{
		Amount:                      float64(amount),
		CardholderPreferredCurrency: true,
		CompletePaymentURL:          successRedirect,
		Country:                     "IN",
		Currency:                    "INR",
		ErrorPaymentURL:             errorRedirect,
		Language:                    "en",
		MerchantReferenceID:         merchantReferenceId,
		Metadata:                    Md{true},
	}

	data, _ := json.Marshal(reqBody)

	req, err := http.NewRequest(http.MethodPost, secret.RAPYD_BASE_URL+"/checkout", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	c := http.DefaultClient
	signer := NewRapydSigner([]byte(secret.RAPYD_ACCESS_KEY), []byte(secret.RAPYD_SECRET_KEY))
	signer.SignRequest(req, data)

	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	d, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(d))
	newData := map[string]interface{}{}
	err = json.Unmarshal([]byte(d), &newData)
	if err != nil {
		return nil, err
	}

	checkoutId := newData["data"].(map[string]interface{})["id"].(string)
	redirectURL := newData["data"].(map[string]interface{})["redirect_url"].(string)
	merchantId := newData["data"].(map[string]interface{})["payment"].(map[string]interface{})["merchant_reference_id"].(string)

	dr := types.DepositResponse{
		PoolId:       poolId,
		UserEmail:    email,
		Amount:       amount,
		CheckoutLink: redirectURL,
		CheckoutId:   checkoutId,
		MerchantId:   merchantId,
	}

	return &dr, nil

}

func RetrieveCheckoutStatus(poolId, email, checkoutId string) (bool, error) {
	// retrieve checkout status
	paid, amount, merchantId, err := getCheckoutStatus(checkoutId)
	if err != nil {
		return false, err
	}
	fmt.Println(paid)
	if paid == false {
		return false, nil
	}

	repo := db.New()

	// check if record exists in transaction db, merchantId acts as an unique key
	doesExists, err := repo.DoesRecordExists(secret.TXN_COLLECTION, merchantId)
	if err != nil {
		return false, nil
	}

	if doesExists {
		return true, nil
	}

	// update the balance in pool, and recalculate pool balance
	u, err := repo.GetRecordDetails(secret.USER_COLLECTIONS, email)
	if err != nil {
		return false, err
	}
	var user types.User
	json.Unmarshal(utils.GetBytesFromInterface(u), &user)
	user.Balance = user.Balance + amount
	um, _ := user.EncodeToMap()
	repo.Save(secret.USER_COLLECTIONS, email, um)

	// update the balance in pool
	p, err := repo.GetPoolDetails(secret.POOL_COLLECTION, poolId)
	if err != nil {
		return false, err
	}

	var pool types.Pool
	json.Unmarshal(utils.GetBytesFromInterface(p), &pool)
	for i := 0; i < len(pool.MembersList); i++ {
		pEmail := pool.MembersList[i].Email
		if pEmail == email {
			pool.MembersList[i].Balance += amount
			break
		}
	}
	// recalculate pool balance an
	totalBalance := float64(0)
	for i := 0; i < len(pool.MembersList); i++ {
		totalBalance += pool.MembersList[i].Balance
	}

	pool.PoolBalance = totalBalance

	pmap, _ := pool.EncodeToMap()
	// update pool
	repo.Save(secret.POOL_COLLECTION, poolId, pmap)

	// also save to transaction db an entry of deposit or withdrawl
	txnLog := types.Transaction{
		PoolId:     poolId,
		UserEmail:  email,
		CheckoutId: checkoutId,
		Amount:     amount,
		Paid:       true,
		TxnType:    types.DEPOSIT,
	}

	txnm, _ := txnLog.EncodeToMap()

	repo.Save(secret.TXN_COLLECTION, merchantId, txnm)
	// TODO::add secure auditing

	return true, nil
}

func getCheckoutStatus(checkoutId string) (bool, float64, string, error) {
	req, err := http.NewRequest(http.MethodGet, secret.RAPYD_BASE_URL+"/checkout/"+checkoutId, nil)
	if err != nil {
		log.Println(err)
		return false, 0.0, "", err
	}

	c := http.DefaultClient
	signer := NewRapydSigner([]byte(secret.RAPYD_ACCESS_KEY), []byte(secret.RAPYD_SECRET_KEY))
	signer.SignRequest(req, nil)

	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return false, 0.0, "", err
	}

	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
	newData := map[string]interface{}{}
	err = json.Unmarshal([]byte(d), &newData)
	if err != nil {
		return false, 0.0, "", err
	}

	paid := newData["data"].(map[string]interface{})["payment"].(map[string]interface{})["paid"].(bool)
	amt := newData["data"].(map[string]interface{})["payment"].(map[string]interface{})["amount"].(float64)
	merchantId := newData["data"].(map[string]interface{})["payment"].(map[string]interface{})["merchant_reference_id"].(string)

	return paid, amt, merchantId, nil
}

func GetSpareLimit(emailId string) (float64, error) {

	repo := db.New()

	isExists, err := repo.DoesRecordExists(secret.USER_COLLECTIONS, emailId)
	if err != nil {
		return -99.0, err
	}

	if !isExists {
		return -99.0, errors.New("User doesn't exists")
	}

	record, err := repo.GetRecordDetails(secret.USER_COLLECTIONS, emailId)
	if err != nil {
		return -99.0, err
	}

	var user types.User
	json.Unmarshal(utils.GetBytesFromInterface(record), &user)

	balance := user.Balance

	// calculate pool balances
	allPools, err := repo.GetAllPools(secret.POOL_COLLECTION)
	if err != nil {
		return -99.0, err
	}

	fmt.Println(allPools)

	ap := *allPools
	uswps := []types.Pool{}

	lockedBalance := float64(0)
	for i := 0; i < len(ap); i++ {
		for j := 0; j < len(ap[i].MembersList); j++ {
			if ap[i].MembersList[j].Email == emailId {
				uswps = append(uswps, ap[i])
				lockedBalance += math.Abs(ap[i].MembersList[j].Balance)
				break
			}
		}
	}

	spareBalance := balance - lockedBalance

	return spareBalance, nil

}

func WithDrawFunds(request types.WithdrawlRequest) (*types.WithdrawlResponse, error) {
	spareLimit, err := GetSpareLimit(request.Email)
	if err != nil {
		return nil, err
	}

	if request.Amount < 100 || request.Amount > spareLimit {
		return nil, errors.New("The amount requested should be greater than equal to 100 and less than spare Balance :: " + fmt.Sprint(spareLimit))
	}

	repo := db.New()

	isExists, err := repo.DoesRecordExists(secret.USER_COLLECTIONS, request.Email)
	if err != nil {
		return nil, err
	}

	if !isExists {
		return nil, errors.New("User doesn't exists in db")
	}

	myUserRecord, err := repo.GetRecordDetails(secret.USER_COLLECTIONS, request.Email)
	if err != nil {
		return nil, err
	}

	var myUser types.User
	json.Unmarshal(utils.GetBytesFromInterface(myUserRecord), &myUser)

	request.FirstName = myUser.FirstName
	request.LastName = myUser.LastName

	// sign a rapyd request for payout and complete the payout
	response, err := CreateAndCompletePayout(request)
	if err != nil {
		return nil, err
	}
	// once payout is completed deduct the balance from user at account level
	x, err := db.New().GetRecordDetails(secret.USER_COLLECTIONS, request.Email)
	if err != nil {
		return nil, err
	}

	var u types.User
	json.Unmarshal(utils.GetBytesFromInterface(x), &u)
	u.Balance -= response.Amount

	sam, _ := u.EncodeToMap()

	// save in db
	repo.Save(secret.USER_COLLECTIONS, u.Email, sam)

	// create a transaction entry
	txnLogs := types.Transaction{
		PayoutId:  response.PayoutId,
		UserEmail: u.Email,
		TxnType:   types.WITHDRAWL,
		Amount:    response.Amount,
		Paid:      true,
	}

	txnm, _ := txnLogs.EncodeToMap()

	repo.Save(secret.TXN_COLLECTION, uuid.New().String(), txnm)

	// TODO: Add a secure auditing logs

	return response, nil
}

func CreateAndCompletePayout(details types.WithdrawlRequest) (*types.WithdrawlResponse, error) {

	type SenderDetails struct {
		CompanyName             string `json:"company_name"`
		IdentificationType      string `json:"identification_type"`
		IdentificationValue     string `json:"identification_value"`
		PhoneNumber             string `json:"phone_number"`
		Occupation              string `json:"occupation"`
		SourceOfIncome          string `json:"source_of_income"`
		DateOfBirth             string `json:"date_of_birth"`
		Address                 string `json:"address"`
		PurposeCode             string `json:"purpose_code"`
		BeneficiaryRelationship string `json:"beneficiary_relationship"`
	}

	type BeneficiaryDetails struct {
		FirstName           string `json:"first_name"`
		LastName            string `json:"last_name"`
		Nationality         string `json:"nationality"`
		IdentificationType  string `json:"identification_type"`
		IdentificationValue string `json:"identification_value"`
		PhoneNumber         string `json:"phone_number"`
		AccountNumber       string `json:"account_number"`
		BankBranchCode      string `json:"bank_branch_code"`
	}

	body := struct {
		PayoutAmount          float64            `json:"payout_amount"`
		ConfirmAutomatically  string             `json:"confirm_automatically"`
		PayoutMethodType      string             `json:"payout_method_type"`
		SenderCurrency        string             `json:"sender_currency"`
		SenderCountry         string             `json:"sender_country"`
		BeneficiaryCountry    string             `json:"beneficiary_country"`
		PayoutCurrency        string             `json:"payout_currency"`
		SenderEntityType      string             `json:"sender_entity_type"`
		BeneficiaryEntityType string             `json:"beneficiary_entity_type"`
		Ewallet               string             `json:"ewallet"`
		Beneficiary           BeneficiaryDetails `json:"beneficiary"`
		Sender                SenderDetails      `json:"sender"`
		Description           string             `json:"description"`
	}{
		PayoutAmount:          details.Amount,
		ConfirmAutomatically:  "true",
		PayoutMethodType:      "in_airtelpaymentsbankltd_bank",
		SenderCurrency:        "USD",
		SenderCountry:         "IN",
		BeneficiaryCountry:    "IN",
		PayoutCurrency:        "INR",
		SenderEntityType:      "company",
		BeneficiaryEntityType: "individual",
		Ewallet:               "ewallet_95ae54567782dcae549bf2836f00c0e4",
		Beneficiary: BeneficiaryDetails{
			FirstName:           details.FirstName,
			LastName:            details.LastName,
			Nationality:         "IN",
			IdentificationType:  "identification_id",
			IdentificationValue: "ZPCUSTOMERTINYFUNDCONNECT11213",
			PhoneNumber:         details.PhoneNumber,
			AccountNumber:       details.BankAccountNumber1,
			BankBranchCode:      details.BankBranchCode,
		},
		Sender: SenderDetails{
			CompanyName:             "Tiny Fund Connect",
			IdentificationType:      "international_passport",
			IdentificationValue:     "Z5111210",
			PhoneNumber:             "0916079100257",
			Occupation:              "FinancialServices",
			SourceOfIncome:          "business_income",
			DateOfBirth:             "07/08/2000",
			Address:                 "Block 4 Koramangla Bengaluru",
			PurposeCode:             "investment_income",
			BeneficiaryRelationship: "customer",
		},
		Description: "Tiny Fund User Balance withdrawl",
	}

	data, _ := json.Marshal(body)

	fmt.Println(fmt.Sprintf("\"%s\"", string(data)))

	req, err := http.NewRequest(http.MethodPost, secret.RAPYD_BASE_URL+"/payouts", bytes.NewBuffer(data))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	c := http.DefaultClient
	signer := NewRapydSigner([]byte(secret.RAPYD_ACCESS_KEY), []byte(secret.RAPYD_SECRET_KEY))
	signer.SignRequest(req, data)

	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
	newData := map[string]interface{}{}
	err = json.Unmarshal([]byte(d), &newData)
	if err != nil {
		return nil, err
	}

	payoutId := newData["data"].(map[string]interface{})["id"].(string)
	amt := newData["data"].(map[string]interface{})["amount"].(float64)
	done, err := completePayout(payoutId)
	if err != nil {
		return nil, err
	}

	if done {
		v := types.WithdrawlResponse{
			Status:   types.SUCCESS_MESSAGE_STATUS_VALUE,
			PayoutId: payoutId,
			Amount:   amt,
		}

		return &v, nil
	}

	return nil, errors.New("Unknown error occured")

}

func completePayout(payoutId string) (bool, error) {
	req, err := http.NewRequest(http.MethodPost, secret.RAPYD_BASE_URL+"/payouts/complete/"+payoutId, nil)
	if err != nil {
		log.Println(err)
		return false, err
	}

	c := http.DefaultClient
	signer := NewRapydSigner([]byte(secret.RAPYD_ACCESS_KEY), []byte(secret.RAPYD_SECRET_KEY))
	signer.SignRequest(req, nil)

	res, err := c.Do(req)
	if err != nil {
		log.Println(err)
		return false, err
	}

	d, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
	newData := map[string]interface{}{}
	err = json.Unmarshal([]byte(d), &newData)
	if err != nil {
		return false, err
	}

	pid := newData["data"].(map[string]interface{})["id"].(string)

	if pid == payoutId {
		return true, nil
	}

	return false, nil
}

func GetAllDepositsByEmail(email string) ([]types.Transaction, error) {
	repo := db.New()

	txnLogs, err := repo.GetAllTxnLogsByEmailId(secret.TXN_COLLECTION, email, types.DEPOSIT)
	if err != nil {
		return nil, err
	}
	return *txnLogs, nil
}

func GetAllWithdrawlsByEmail(email string) ([]types.Transaction, error) {
	repo := db.New()

	txnLogs, err := repo.GetAllTxnLogsByEmailId(secret.TXN_COLLECTION, email, types.WITHDRAWL)
	if err != nil {
		return nil, err
	}
	return *txnLogs, nil
}

func GetAllDepositsByPoolId(poolId string) ([]types.Transaction, error) {
	repo := db.New()

	txnLogs, err := repo.GetAllTxnLogsByPoolId(secret.TXN_COLLECTION, poolId, types.DEPOSIT)
	if err != nil {
		return nil, err
	}
	return *txnLogs, nil
}
