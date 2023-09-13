<p align="center"><img src="https://pin.ski/3smqXEU" alt="Brand logo of tiny fund connect"></p>

# TinyFundConnect
In today’s financial landscape, several critical issues hold back potential investors, limiting their access to the cryptocurrency market. Limited financial resources, Information asymmetry and Lack of Solution are few factors. Tiny Fund connect fills this market void and brings to you the next generation Peer to peer Micro investing application, where users can form pools, join pools, pool in their financial resources, trade cryptocurrency and share profit collectively.

# Try it out
[try it out here](http://172.232.132.251/)

# Installation

- Download and Install golang https://go.dev/doc/install
- clone the repository using
  ```
  git clone https://github.com/sap200/TinyFundConnect.git
  ```
- Open a binance spot test net account https://testnet.binance.vision/
- open a pangea account https://pangea.cloud/
- open a rapyd account https://www.rapyd.net/
- Fill in this values which are descriptive enough
```
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
```
- go inside the directory and enter these commands
  ```
  go build -o tiny_fund_connect
  ./tiny_fund_connect
  ```

# Use of Pangea SDK

- AuthN was used for signup, login and password reset flow
- Embargo, user_intel and ip_intel was used to validate the incoming request source
- Redact was used while pool chatting and poll creation to ensure zero tolerance against profanity and sharing of personal details
- secure auditing was used at Login, Signup, Password reset, Pool Creation, Order creation, Deposits and Withdrawls.
- The directories with pangea is pangeaauth, pangeachecks and pangealogger, the integration can be seen in authhandler, poolhandler, tradehandler, paymenthandler
  
  
