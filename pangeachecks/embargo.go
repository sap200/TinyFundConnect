package pangeachecks

import (
	"context"
	"encoding/json"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/embargo"

	"github.com/sap200/TinyFundConnect/secret"
)

var checkedIPList = make(map[string]string)

const (
	ALLOW_ACCESS      = "allow"
	DONT_ALLOW_ACCESS = "dont_allow"
)

func IsIpInSanctionedList(ip string) (bool, error) {

	val, ok := checkedIPList[ip]
	if ok {
		if val == ALLOW_ACCESS {
			return false, nil
		} else {
			return true, nil
		}
	}

	token := secret.PANGEA_EMBARGO_TOKEN

	embargocli := embargo.New(&pangea.Config{
		Token:  token,
		Domain: secret.PANGEA_DOMAIN,
	})

	ctx := context.Background()
	input := &embargo.IPCheckRequest{
		IP: ip,
	}

	checkResponse, err := embargocli.IPCheck(ctx, input)
	if err != nil {
		return false, err
	}

	var c embargo.CheckResult
	json.Unmarshal([]byte(pangea.Stringify(checkResponse.Result)), &c)

	if c.Count == 0 {
		checkedIPList[ip] = ALLOW_ACCESS
	} else {
		checkedIPList[ip] = DONT_ALLOW_ACCESS

	}

	return c.Count != 0, nil
}
