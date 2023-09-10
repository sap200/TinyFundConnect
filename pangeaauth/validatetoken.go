package pangeaauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/sap200/TinyFundConnect/secret"
)

func ValidatePangeaToken(token string) error {
	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.ClientTokenCheckRequest{
		Token: token,
	}

	resp, err := authncli.Client.Token.Check(ctx, input)
	if err != nil {
		return err
	}

	fmt.Println("Response:: ", pangea.Stringify(resp.Result))
	// add a line that response should contain

	return nil

}

func ValidateEmailPangeaToken(email, token string) error {

	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.ClientTokenCheckRequest{
		Token: token,
	}

	resp, err := authncli.Client.Token.Check(ctx, input)
	if err != nil {
		return err
	}

	fmt.Println("Response:: ", pangea.Stringify(resp.Result))

	var tr authn.ClientTokenCheckResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &tr)

	if tr.Email != email {
		return errors.New("invalid email for token")
	}

	// add a line that response should contain

	return nil
}
