package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/db"
	e "github.com/sap200/TinyFundConnect/errors"
	"github.com/sap200/TinyFundConnect/pangeaauth"
	"github.com/sap200/TinyFundConnect/pangeachecks"
	"github.com/sap200/TinyFundConnect/pangealogger"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
	"github.com/sap200/TinyFundConnect/utils"
)

func SignUpHandler(c *gin.Context) {

	// get 2 parameters
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// validate Password
	passResult := utils.ValidatePassword(user.Password)
	if !passResult {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_PASSWORD_FORMAT,
		})
		return
	}

	repo := db.New()

	ok, err := repo.DoesRecordExists(secret.USER_COLLECTIONS, user.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.EMAIL_ALREADY_EXISTS_IN_DB,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}
	ok2, err := pangeachecks.IsUserEmailInBreachedList(user.Email)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
				"message": err.Error(),
			})
			return
		}
	}
	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok2 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok2 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_EMAIL_BREACHED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": errorMessage,
		})
		return
	}

	updatedUser, err := pangeaauth.InitiateSignupFlow(user.Email, user.Password, user.FirstName, user.LastName, user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          types.SUCCESS_MESSAGE_CODE,
		"message":       types.SUCCESS_MESSAGE_STATUS_VALUE,
		"email":         updatedUser.Email,
		"firstName":     updatedUser.FirstName,
		"lastName":      updatedUser.LastName,
		"emailVerified": updatedUser.EmailVerified,
		"customerId":    updatedUser.CustomerId,
		"flowId":        updatedUser.FlowId,
	})
}

func ResendEmailVerificationMailHandler(c *gin.Context) {
	// get 2 parameters
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}
	ok2, err := pangeachecks.IsUserEmailInBreachedList(user.Email)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
				"message": err.Error(),
			})
			return
		}
	}
	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok2 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok2 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_EMAIL_BREACHED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": errorMessage,
		})
		return
	}

	// check if email is already verified do not verify again and send proper code
	isVerified, err := pangeaauth.GetEmailVerificationStatus(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if isVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.EMAIL_ALREADY_VERIFIED,
			"message": e.EMAIL_AREADY_VERIFIED_MESSAGE,
		})

		return
	}

	err = pangeaauth.ResendEmailVerification(user.Email, user.FlowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    types.SUCCESS_MESSAGE_CODE,
		"message": types.SUCCESS_MESSAGE_STATUS_VALUE,
	})

}

func CheckIfEmailIsVerified(c *gin.Context) {
	// since this is a get api and no db writing is there
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": errorMessage,
		})
		return
	}

	isVerified, err := pangeaauth.GetEmailVerificationStatus(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// update db record from our side if both entries are not same
	repo := db.New()
	userDetails, _ := repo.GetRecordDetails(secret.USER_COLLECTIONS, user.Email)
	userDet, _ := userDetails.(types.User)
	if userDet.EmailVerified != isVerified {
		repo.MarkUserVerified(secret.USER_COLLECTIONS, user.Email, "emailVerified", isVerified)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       types.SUCCESS_MESSAGE_CODE,
		"message":    types.SUCCESS_MESSAGE_STATUS_VALUE,
		"email":      user.Email,
		"isVerified": isVerified,
	})
}

func SendEmailOtpDuringSignup(c *gin.Context) {

	// since this is a get api and no db writing is there
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}
	ok2, err := pangeachecks.IsUserEmailInBreachedList(user.Email)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
				"message": err.Error(),
			})
			return
		}
	}
	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok2 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok2 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_EMAIL_BREACHED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": errorMessage,
		})
		return
	}

	// everything looks alright, call the enroll MFA
	err = pangeaauth.EnrollMFA(user.Email, user.FlowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return

	}
	fmt.Println(err)

	c.JSON(http.StatusOK, gin.H{
		"code":    types.SUCCESS_MESSAGE_STATUS_VALUE,
		"message": "Otp Sent to your email",
	})
}

func VerifyEmailOtpDuringSignup(c *gin.Context) {
	// since this is a get api and no db writing is there
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}
	ok2, err := pangeachecks.IsUserEmailInBreachedList(user.Email)
	if err != nil {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
				"message": err.Error(),
			})
			return
		}
	}
	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok2 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok2 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_EMAIL_BREACHED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": errorMessage,
		})
		return
	}

	// call the function
	err = pangeaauth.CompleteEnrollMFAComplete(user.Email, user.FlowId, user.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// complete the signup or signin flow to receive the tokens

	fresp, err := pangeaauth.CompleteSignupOrLoginFlow(user.FlowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	// log when user has completed signup
	event := pangealogger.New(
		"User with email "+user.Email+" has signed up for Tiny Fund Connect",
		"SIGNUP_EVENT",
		user.Ip,
		user.Email,
		"",
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, gin.H{
		"status":             types.SUCCESS_MESSAGE_STATUS_VALUE,
		"code":               types.SUCCESS_MESSAGE_CODE,
		"redirectURL":        types.REDIRECT_URL_AFTER_SUCCESSFUL_SIGNUP,
		"activeToken":        fresp.ActiveToken.Token,
		"activeTokenId":      fresp.ActiveToken.ID,
		"activeTokenExpiry":  fresp.ActiveToken.Expire,
		"refreshToken":       fresp.ActiveToken.Token,
		"refreshTokenId":     fresp.ActiveToken.ID,
		"refreshTokenExpiry": fresp.RefreshToken.Expire,
	})
}

func ValidateTokenStatus(c *gin.Context) {
	// since this is a get api and no db writing is there
	var tokens types.Tokens
	err := c.ShouldBindJSON(&tokens)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	err = pangeaauth.ValidatePangeaToken(tokens.ActiveToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.BAD_REQUEST_STATUS_CODE_SIGNUP,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  types.SUCCESS_MESSAGE_STATUS_VALUE,
		"code":    types.SUCCESS_MESSAGE_CODE,
		"message": "Token is a valid.",
		"isValid": true,
	})

}

func StartLogin(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// verify email is verified
	isEmailVerified, err := pangeaauth.GetEmailVerificationStatus(user.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": "Error during email verification",
		})
		return
	}

	if !isEmailVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": "Please verify your email before logging in!",
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": errorMessage,
		})
		return
	}

	u, err := pangeaauth.StartLogin(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    types.LOGIN_SUCCESS_CODE,
		"message": types.LOGIN_SUCCESS_MFA_START_MESSAGE,
		"flowId":  u.FlowId,
	})
}

func CompleteLoginHandler(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": errorMessage,
		})
		return
	}

	fresp, err := pangeaauth.CompleteLogin(user.FlowId, user.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	// log when user has completed signup
	event := pangealogger.New(
		"User with email "+user.Email+" has Logged in",
		"LOGIN_EVENT",
		user.Ip,
		user.Email,
		"",
	)
	pangealogger.Log(event)

	c.JSON(http.StatusOK, gin.H{
		"status":             types.SUCCESS_MESSAGE_STATUS_VALUE,
		"code":               types.LOGIN_SUCCESS_CODE,
		"redirectURL":        types.REDIRECT_URL_AFTER_LOGIN,
		"activeToken":        fresp.ActiveToken.Token,
		"activeTokenId":      fresp.ActiveToken.ID,
		"activeTokenExpiry":  fresp.ActiveToken.Expire,
		"refreshToken":       fresp.ActiveToken.Token,
		"refreshTokenId":     fresp.ActiveToken.ID,
		"refreshTokenExpiry": fresp.RefreshToken.Expire,
	})
}

func LogoutHandler(c *gin.Context) {
	// since this is a get api and no db writing is there
	var tokens types.Tokens
	err := c.ShouldBindJSON(&tokens)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGOUT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	s, err := pangeaauth.LogoutFromSession(tokens.ActiveToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGOUT_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if s != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    types.LOGOUT_SUCCESS_CODE,
			"message": types.LOGOUT_SUCCESS_MESSAGE,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"code":    e.LOGOUT_FAILURE_CODE,
		"message": err.Error(),
	})
	return

}

func PasswordResetEmailVerifyHandler(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	// validate email
	res := utils.ValidateEmail(user.Email)
	if !res {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": e.INVALID_EMAIL_FORMAT,
		})
		return
	}

	// pangea checks on IP and email
	ok1, err := pangeachecks.IsIpInSanctionedList(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	ok3, err := pangeachecks.IsTheIpVpnOrProxyOrMalicious(user.Ip)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if ok1 || ok3 {
		errorMessage := ""
		if ok1 {
			errorMessage += e.PANGEA_CHECK_IP_SANCTIONED
		}

		if ok3 {
			if errorMessage != "" {
				errorMessage += " | "
			}
			errorMessage += e.PANGEA_CHECK_IP_INTEL
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": errorMessage,
		})
		return
	}

	flowId, err := pangeaauth.VerifyEmailForPasswordReset(user.Email)
	if flowId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": "Flow Id is an empty string",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": types.SUCCESS_MESSAGE_STATUS_VALUE,
		"code":   types.LOGIN_SUCCESS_CODE,
		"flowId": flowId,
	})

}

func PasswordResetMFAVerificationHandler(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}
	fid, err := pangeaauth.EnterPasswordAndStartMFAVerification(user.FlowId, user.Password, user.Cancel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if fid == types.PASSWORD_RESET_CANCELLATION_MESSAGE {
		c.JSON(http.StatusOK, gin.H{
			"code":    types.SUCCESS_MESSAGE_CODE,
			"message": "Password reset cancelled",
		})
		return
	}

	if fid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": "FlowId is empty, please verify email",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":   types.LOGIN_SUCCESS_CODE,
		"status": types.SUCCESS_MESSAGE_STATUS_VALUE,
		"flowId": fid,
	})
}

func PasswordResetMfaCompleteHandler(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	x, err := pangeaauth.ConfirmPasswordChangeWithMFAVerification(user.FlowId, user.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	if x {
		c.JSON(http.StatusOK, gin.H{
			"status":  types.SUCCESS_MESSAGE_STATUS_VALUE,
			"code":    types.LOGIN_SUCCESS_CODE,
			"message": "password reset success",
		})
	}
}

func CheckIfPasswordResetEmailVerified(c *gin.Context) {
	var user types.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	x, err := pangeaauth.ShouldWeCircleInVerifyEmailPage(user.FlowId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    e.LOGIN_FAILURE_CODE,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verified": !x,
	})

}
