package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sap200/TinyFundConnect/handler"
	"github.com/sap200/TinyFundConnect/types"
)

func New() *gin.Engine {
	r := gin.Default()
	// Enable CORS for all origins
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // Allow any origin
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, token")
		c.Header("Access-Control-Allow-Credentials", "true") // Optional: Set this header if you need to include credentials (e.g., cookies) in your requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // Handle preflight (OPTIONS) requests
			return
		}
		c.Next()
	})
	return r
}

func RegisterPaths(r *gin.Engine) {
	r.GET(types.TEST_PATH, handler.TestPathHandler)
	r.POST(types.SIGNUP_PATH, handler.SignUpHandler)
	r.POST(types.VERIFY_RESEND_EMAIL_LINK, handler.ResendEmailVerificationMailHandler)
	r.POST(types.CHECK_IS_EMAIL_VERIFIED, handler.CheckIfEmailIsVerified)
	r.POST(types.ENROLL_MFA, handler.SendEmailOtpDuringSignup)
	r.POST(types.ENROLL_MFA_COMPLETE, handler.VerifyEmailOtpDuringSignup)
	r.POST(types.VALIDATE_TOKEN_STATUS, handler.ValidateTokenStatus)
	r.POST(types.START_LOGIN, handler.StartLogin)
	r.POST(types.COMPLETE_LOGIN, handler.CompleteLoginHandler)
	r.POST(types.LOGOUT, handler.LogoutHandler)
	r.POST(types.PASSWORD_RESET_VERIFY_EMAIL, handler.PasswordResetEmailVerifyHandler)
	r.POST(types.PASSWORD_RESET_ENTER_PASSWORD_AND_START_MFA_VERIFICATION, handler.PasswordResetMFAVerificationHandler)
	r.POST(types.PASSWORD_RESET_COMPLETE_MFA, handler.PasswordResetMfaCompleteHandler)
	r.POST(types.PASSWORD_RESET_CHECK_EMAIL_VERIFIED, handler.CheckIfPasswordResetEmailVerified)
	// Pool Creation
	r.POST(types.CREATE_POOL, handler.CreatePoolHandler)
	r.GET(types.GET_POOL, handler.GetPoolDetailsHandler)
	r.PUT(types.POOL_MEMBER_STATUS_UPDATE, handler.UpdatePoolMemberStatus)
	r.PUT(types.POOL_ADD_MEMBER, handler.PoolAddMemberHandler)
	r.PUT(types.POOL_EXIT, handler.ExitFromPoolHandler)
	r.GET(types.GET_ALL_POOLS_WITH_EMAIL, handler.GetAllPoolsWithListedEmail)
	//payments
	r.POST(types.DEPOSIT_START, handler.StartDepositHandler)
	r.POST(types.DEPOSIT_STATUS_CHECK, handler.DepositStatusCheckHandler)
	r.GET(types.GET_WITHDRAWL_LIMIT, handler.GetWithdrawlLimitHandler)
	r.POST(types.WITHDRAW_START, handler.WithdrawlStartHandler)
	r.GET(types.WITHDRAW_GETALL, handler.WithdrawlGetAllHandler)
	r.GET(types.DEPOSITS_GETALL, handler.DepositGetAllHandler)
	r.GET(types.DEPOSIT_GETALL_BY_POOL, handler.DepositGetAllByPoolHandler)
	// chat apis
	r.POST(types.SEND_CHAT, handler.SendChatHandler)
	r.GET(types.GET_CHAT_BY_POOL, handler.GetChatByPoolHandler)
	// Poll
	r.POST(types.CREATE_POLL, handler.CreatePollHandler)
	r.PUT(types.INACTIVATE_POLL, handler.InactivatePollHandler)
	r.POST(types.CREATE_VOTE, handler.CreateVoteHandler)
	r.GET(types.GET_POLL_BY_ID, handler.GetAllPollById)
	r.GET(types.GET_ALL_POLL_BY_POOL, handler.GetAllPollByPoolHandler)
	r.GET(types.USER_ALREADY_VOTED_FOR_POLL, handler.CheckIfUserHasAlreadyVotedForPoll)
	r.GET(types.GET_POLL_RESULT, handler.GetPollResultHandler)
	// TRADE apis
	r.GET(types.GET_CANDLE_DATA, handler.GetCandleDatHandler)
	r.GET(types.GET_ORDER_BOOK, handler.RetrieveOrderBookHandler)
	r.GET(types.GET_RECENT_TRADE, handler.GetRecentTradeHandler)
	r.GET(types.GET_ORDERS_BY_POOL_ID, handler.GetTradeOrderByPoolIdHandler)
	r.GET(types.GET_OPEN_ORDERS_BY_POOL_ID, handler.GetTradeOpenOrdersByPoolId)
	r.DELETE(types.CANCEL_OPEN_ORDER_BY_ORDER_ID, handler.CancelOpenOrderTradeHandler)
	r.POST(types.CREATE_A_LIMIT_ORDER, handler.CreateLimitOrderTradeHandler)
	r.GET(types.GET_ALL_ALLOWED_SYMBOLS, handler.GetAllSymbolTradeHandler)

	// UTILITY
	r.GET(types.CHECK_TOKEN, handler.CheckTokenHandler)
}
