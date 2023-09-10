package types

const (
	SUCCESS_MESSAGE_STATUS_VALUE                            = "SUCCESS"
	SUCCESS_MESSAGE_CODE                                    = "TFC000"
	SUCCESS_MESSAGE_STATUS_KEY                              = "status"
	SUCCESS_MESSAGE_STATUS_CODE                             = 200
	ERROR_MESSAGE_STATUS_VALUE                              = "FAILURE"
	ERROR_MESSAGE_STATUS_KEY                                = "message"
	ERROR_MESSAGE_STATUS_CODE                               = 400
	REDIRECT_URL_AFTER_SUCCESSFUL_SIGNUP                    = "https://localhost:8080/login"
	REDIRECT_URL_AFTER_LOGIN                                = "http://localhost:8080/dashboard"
	MFA_VERIFICATION_PATH                                   = "verify/mfa/start"
	MFA_VERIFICATION_COMPLETE_PATH                          = "verify/mfa/complete"
	LOGIN_SUCCESS_CODE                                      = "LOGIN000"
	LOGIN_SUCCESS_MFA_START_MESSAGE                         = "Please verify your email otp"
	LOGIN_COMPLETE_PATH                                     = "complete"
	LOGOUT_SUCCESS_MESSAGE                                  = "Logged out the current user login session"
	LOGOUT_SUCCESS_CODE                                     = "LOGOUT000"
	PANGEA_MAIL_VERIFY_PASSWORD_RESET_FINAL_PATH            = "verify/password_reset"
	PANGEA_MOVE_FROM_EMAIL_VERIFICATION_PATH_PASSWORD_RESET = "reset/password"
	PANGEA_PASSWORD_RESET_CANCELLED_NEXT_PATH               = "verify/password"
	PASSWORD_RESET_CANCELLATION_MESSAGE                     = "Password Reset Request Cancelled Successfully"
	TOKEN_VAR                                               = "token"
)
