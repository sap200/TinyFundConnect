package secret

const (
	PROJECT_ID                 = "tinyfundconnect"
	FROM_EMAIL                 = "<<your-business-email-id>>"
	FROM_EMAIL_PASSWORD        = "<<email-password>>"
	HOST                       = "<<email-host>>"
	EMAIL_PATH_STRING          = ":email"
	UPDATE_VERIFIED_FIELD_NAME = "emailVerified"

	// Collections name
	USER_COLLECTIONS     = "user_records"
	USER_COLLECTIONS_KEY = "email"
	POOL_COLLECTION      = "pool_records"
	TXN_COLLECTION       = "transaction_records"
	CHAT_COLLECTION      = "chat_records"
	POLLING_COLLECTION   = "polling_records"
	ORDERS_COLLECTION    = "closed_order_records"

	// pangea
	PANGEA_EMBARGO_TOKEN          = "<<pangea-embargo-token>>"
	PANGEA_DOMAIN                 = "<<pangea-domain>>"
	PANGEA_INTEL_TOKEN            = "<<pangea-intel-token>>"
	PANGEA_AUTHN_TOKEN            = "<<pangea-auth-token>>"
	PANGEA_AUDIT_SCHEMA_CONFIG_ID = "<<pangea-audit-config-schema>>"

	// auth
	PANGEA_SIGNUP_CALLBACK_URL = "<<pangea-callback-url-on-email-verification>>"

	// timeout
	CANCEL_TIMEOUT_PANGEA = 60

	// RAPYD URLS
	RAPYD_BASE_URL   = "<<rapyd-base-url>>"
	RAPYD_ACCESS_KEY = "<<rapyd-access-key>>"
	RAPYD_SECRET_KEY = "<<rapyd-secret-key>>"

	// TRADING
	BINANCE_API_KEY                  = "<<binance-api-key>>"
	BINANCE_SECRET                   = "<<binance-secret-key>>"
	BINANCE_SPOT_TEST_NET_API_KEY    = "<<binance-spot-test-net-api-key>>"
	BINANCE_SPOT_TEST_NET_SECRET_KEY = "<<binance-spot-test-net-secret-key>>"
	BINANCE_SPOT_TEST_NET_ENDPOINT   = "<<binance-api-host>>"
)
