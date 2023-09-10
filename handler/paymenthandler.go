package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/db"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangeaauth"
	"github.com/sap200/TinyFundConnect/pangealogger"
	"github.com/sap200/TinyFundConnect/payments"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"github.com/sap200/TinyFundConnect/utils"
)

func StartDepositHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var depReq types.DepositRequest
	err = c.ShouldBindJSON(&depReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	repo := db.New()
	ok, err := repo.DoesRecordExists(secret.POOL_COLLECTION, depReq.PoolId)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	fmt.Println("Reached Here : 001")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.POOL_FAILURE_CODE,
			"message": "Pool with given Id not found",
		})

		return
	}

	fmt.Println("Reached Here : 1")

	ok, err = repo.DoesRecordExists(secret.USER_COLLECTIONS, depReq.UserEmail)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}
	fmt.Println("Reached Here : 2")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": "Email doesn't exists in db",
		})
		return
	}

	fmt.Println("Reached Here : 3")

	depRes, err := payments.StartDeposit(depReq.PoolId, depReq.UserEmail, depReq.SuccessRedirectLink, depReq.ErrorRedirectLink, depReq.Amount)
	fmt.Println(depRes)
	fmt.Println(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	fmt.Println("Reached Here 4")

	nm, _ := depRes.EncodeToMap()

	fmt.Println("Reached Here: 5")

	c.JSON(http.StatusOK, nm)
}

// Deposit status check
func DepositStatusCheckHandler(c *gin.Context) {

	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var depRes types.DepositResponse
	err = c.ShouldBindJSON(&depRes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	repo := db.New()

	// check if email exists

	ok, err := repo.DoesRecordExists(secret.USER_COLLECTIONS, depRes.UserEmail)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": "Email id doesn't exists in user db",
		})

		return
	}

	// check if pool exists
	ok, err = repo.DoesRecordExists(secret.POOL_COLLECTION, depRes.PoolId)
	fmt.Println(ok, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": "Pool Id doesn't exists in pool db",
		})

		return
	}

	// check if email exists in pool
	if ok {
		p, err := repo.GetPoolDetails(secret.POOL_COLLECTION, depRes.PoolId)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.PAYMENT_FAILURE_CODE,
				"message": err.Error(),
			})

			return
		}

		var pool types.Pool
		json.Unmarshal(utils.GetBytesFromInterface(p), &pool)
		userExistsInPool := false
		for i := 0; i < len(pool.MembersList); i++ {
			if pool.MembersList[i].Email == depRes.UserEmail {
				userExistsInPool = true
				break
			}
		}

		if !userExistsInPool {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.PAYMENT_FAILURE_CODE,
				"message": "User not found in pool",
			})

			return
		}
	}

	paid, err := payments.RetrieveCheckoutStatus(depRes.PoolId, depRes.UserEmail, depRes.CheckoutId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})

		return
	}

	if !paid {
		c.JSON(http.StatusOK, gin.H{
			"paid": false,
		})

		return
	}

	// log when user has completed signup
	event := pangealogger.New(
		"User with email "+depRes.UserEmail+" has deposited an amount in pool id  "+depRes.PoolId,
		"DEPOSIT_EVENT",
		"",
		depRes.UserEmail,
		depRes.PoolId,
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, gin.H{
		"paid": true,
	})

}

func GetWithdrawlLimitHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	email := c.Query("email")

	amnt, err := payments.GetSpareLimit(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   e.PAYMENT_SUCCESS_CODE,
		"status": types.SUCCESS_MESSAGE_CODE,
		"amount": amnt,
	})

}

func WithdrawlStartHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	var request types.WithdrawlRequest
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if request.BankAccountNumber1 != request.BankAccountNumber2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": "Bank account number mismatch",
		})
		return
	}

	withdrawlResponse, err := payments.WithDrawFunds(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	event := pangealogger.New(
		"User with email "+request.Email+" withdrew an amount of "+fmt.Sprintf("₹ %.2f", request.Amount),
		"WITHDRAWL_EVENT",
		"",
		request.Email,
		"",
	)
	pangealogger.Log(event)

	rm, _ := withdrawlResponse.EncodeToMap()

	c.JSON(http.StatusOK, rm)

}

func WithdrawlGetAllHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	email := c.Query("email")

	tlog, err := payments.GetAllWithdrawlsByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tlog)
}

func DepositGetAllHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	email := c.Query("email")

	tlog, err := payments.GetAllDepositsByEmail(email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tlog)

}

func DepositGetAllByPoolHandler(c *gin.Context) {
	authToken := c.GetHeader(types.TOKEN_VAR)
	err := pangeaauth.ValidatePangeaToken(authToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    e.UNAUTHORIZED,
			"message": err.Error(),
		})
		return
	}

	poolId := c.Query("poolId")

	tlog, err := payments.GetAllDepositsByPoolId(poolId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.PAYMENT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tlog)

}
