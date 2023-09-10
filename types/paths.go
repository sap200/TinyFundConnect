package types

const (
	TEST_PATH = "/test"

	// SIGNUP
	SIGNUP_PATH              = "/signup"
	VERIFY_EMAIL_PATH        = "/verify/:email"
	VERIFY_RESEND_EMAIL_LINK = "/emailverification/resend"
	CHECK_IS_EMAIL_VERIFIED  = "/check/emailverified"
	ENROLL_MFA               = "/enrollmfa/emailotp"
	ENROLL_MFA_COMPLETE      = "/enrollmfa/verifyemailotp"

	// VALIDATE_TOKEN
	VALIDATE_TOKEN_STATUS = "/validate/token"

	// LOGIN
	START_LOGIN    = "/login/start"
	COMPLETE_LOGIN = "/login/complete"

	// LOGOUT
	LOGOUT = "/logout"

	// PASSWORD RESET
	PASSWORD_RESET_VERIFY_EMAIL                              = "/password/verifyemail"
	PASSWORD_RESET_ENTER_PASSWORD_AND_START_MFA_VERIFICATION = "/password/startmfa"
	PASSWORD_RESET_COMPLETE_MFA                              = "/password/completemfa"
	PASSWORD_RESET_CHECK_EMAIL_VERIFIED                      = "/password/checkemailverified"

	// POOL
	CREATE_POOL               = "/pool/create"
	GET_POOL                  = "/pool/details"
	POOL_MEMBER_STATUS_UPDATE = "/pool/statusupdate"
	POOL_ADD_MEMBER           = "/pool/addmember"
	POOL_EXIT                 = "/pool/exit"
	GET_ALL_POOLS_WITH_EMAIL  = "/pool/getall"

	// payments
	DEPOSIT_START          = "/deposit/start"
	DEPOSIT_STATUS_CHECK   = "/deposit/status"
	GET_WITHDRAWL_LIMIT    = "/withdraw/limit"
	WITHDRAW_START         = "/withdraw/start"
	DEPOSITS_GETALL        = "/deposit/getall"
	WITHDRAW_GETALL        = "/withdraw/getall"
	DEPOSIT_GETALL_BY_POOL = "/deposit/getallbypoolid"

	// chat
	SEND_CHAT        = "/chat/send"
	GET_CHAT_BY_POOL = "/chat/getallbypool"

	// Poll
	CREATE_POLL                 = "/poll/create"
	INACTIVATE_POLL             = "/poll/inactivate"
	CREATE_VOTE                 = "/vote/create"
	GET_ALL_POLL_BY_POOL        = "/poll/getall"
	GET_POLL_BY_ID              = "/poll/get"
	USER_ALREADY_VOTED_FOR_POLL = "/poll/alreadyvoted"
	GET_POLL_RESULT             = "/poll/getresult"

	// TRADE
	GET_CANDLE_DATA               = "/candle/data"
	GET_ORDER_BOOK                = "/orderbook/get"
	GET_RECENT_TRADE              = "/traderecent/get"
	GET_ORDERS_BY_POOL_ID         = "/orders/getbypoolid"
	GET_OPEN_ORDERS_BY_POOL_ID    = "/orders/getopenordersbypoolid"
	CANCEL_OPEN_ORDER_BY_ORDER_ID = "/orders/cancel"
	CREATE_A_LIMIT_ORDER          = "/orders/create"
	GET_ALL_ALLOWED_SYMBOLS       = "/symbols/get"

	// UTILITY API
	CHECK_TOKEN = "/token/check"
)
