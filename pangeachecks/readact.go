package pangeachecks

import (
	"context"
	"encoding/json"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/redact"
	"github.com/sap200/TinyFundConnect/secret"
)

func RedactAMessage(message string) (string, error) {

	token := secret.PANGEA_INTEL_TOKEN

	redactcli := redact.New(&pangea.Config{
		Token:  token,
		Domain: secret.PANGEA_DOMAIN,
	})

	ctx := context.Background()
	input := &redact.TextRequest{
		Text: pangea.String(message),
	}

	redactResponse, err := redactcli.Redact(ctx, input)
	if err != nil {
		return "", err
	}

	var result redact.TextResult
	json.Unmarshal([]byte(pangea.Stringify(redactResponse.Result)), &result)

	return *result.RedactedText, nil
}
