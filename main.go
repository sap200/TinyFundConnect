package main

import (
	"github.com/sap200/TinyFundConnect/router"
)

func main() {
	//test()
	// main function
	r := router.New()
	router.RegisterPaths(r)
	r.Run()
}

/**
* This is a test function to test individual api functionality
 */
func test() {
	// x := pangeaauth.InitiateSignupFlow("saptarsihalder29@gmail.com", "Saptarsi@123", "saptarsi", "halder", "12.12.12")
	// fmt.Println(x)
	//fmt.Println(pangeaauth.ResendEmailVerification("saptarsihalder29@gmail.com", "pfl_gyxemuekvodzzvxkuh7yic2j3nfjpilq"))
	//fmt.Println(pangeaauth.ValidatePangeaToken("abcd"))
	//fmt.Println(pangeaauth.StartLogin("www.saptarsi@gmail.com", "Saptarsi@123"))
	//fmt.Println(pangeaauth.LogoutFromSession("ptu_ni5kfhov76uhtb6ar6ubis66ymcx54aa"))
	//fmt.Println(pangeaauth.ShouldWeCircleInVerifyEmailPage("pfl_63am4sexrdbfsrt7bgb6x2g4fxmlt7pj"))
	// fmt.Println(pangeaauth.ValidatePangeaToken("ptu_g3l7ijs3qqlq5embefoldb44x7qwbmkr"))
	// fmt.Println(payments.StartDeposit("abcd", "www.saptarsi@gmail.com", "http://localhost:8080/test", "http://localhost:8080/test", 12345.00))
	//fmt.Println(payments.RetrieveCheckoutStatus("abc", "def", "checkout_f078aed9f98030b6afbdcc6a098c565f"))
	// fmt.Println(payments.WithDrawFunds(types.WithdrawlRequest{
	// 	FirstName:          "abc",
	// 	LastName:           "def",
	// 	BankAccountNumber1: "01211",
	// 	BankAccountNumber2: "01211",
	// 	BankBranchCode:     "ABC00120",
	// 	PhoneNumber:        "0918900933078",
	// 	Email:              "saptarsihalder29@gmail.com",
	// 	Amount:             151.123,
	// }))

	// fmt.Println(binance.GetMarketData(map[string]string{
	// 	"symbol":   "123",
	// 	"interval": "5m",
	// }))

	// fmt.Println(binance.GetOrderBook(map[string]string{
	// 	"symbol": "BNBUSDT",
	// 	"limit":  "10",
	// }))

	// fmt.Println(pangeachecks.RedactAMessage("My phone number is +916379600458"))

	// e := pangealogger.New("Hello I am so and so here", "ABCD_EVENT", "", "www.saptarsi@gmail.com", "abcd-efgh-123-000")

	// pangealogger.Log(e)

	// //main function
	// go binance.ApportionProfitBetweenPoolMembers(-100, "95c21110-9552-44ea-b285-be043c3c3497")
	// go binance.SaveEarnedFromOrder(60.0, "112136", "pool-01")
	// go func() {
	// 	fmt.Println(binance.GetTotalEarnedFromOrderId("112136", "pool-01"))
	// }()
}
