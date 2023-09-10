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
	"github.com/sap200/TinyFundConnect/types"
)

func StartLogin(emailId, password string) (*types.User, error) {
	// start the flow
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()
	fmt.Println("Here 1")
	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	fmt.Println("here2")
	fts := []authn.FlowType{authn.FTsignin}
	// start flow
	var provider authn.IDProvider
	provider = authn.IDPPassword
	input := authn.FlowStartRequest{
		CBURI:     secret.PANGEA_SIGNUP_CALLBACK_URL,
		Email:     emailId,
		FlowTypes: fts,
		Provider:  &provider,
	}
	fmt.Println("here3")

	resp, err := authncli.Flow.Start(ctx, input)
	if err != nil {
		return nil, err
	}
	var fr authn.FlowStartResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &fr)
	flowId := fr.FlowID
	if flowId == "" {
		return nil, errors.New("Invalid Flow Id: " + pangea.Stringify(resp.Result))
	}
	fmt.Println("FLow Id here: ", flowId)

	// start login flow
	passwordVerificationRequest := authn.FlowVerifyPasswordRequest{
		FlowID:   flowId,
		Password: pangea.String(password),
	}

	passVerRes, err := authncli.Flow.Verify.Password(ctx, passwordVerificationRequest)
	if err != nil {
		return nil, err
	}

	var pr authn.FlowVerifyPasswordResult
	json.Unmarshal([]byte(pangea.Stringify(passVerRes.Result)), &pr)

	if pr.Error != nil {
		return nil, errors.New(*pr.Error)
	}

	fmt.Println(pr)

	if pr.NextStep == types.MFA_VERIFICATION_PATH {
		mfaVerReq := authn.FlowVerifyMFAStartRequest{
			FlowID:      flowId,
			MFAProvider: authn.MFAPEmailOTP,
		}

		mfaVerResp, err := authncli.Flow.Verify.MFA.Start(ctx, mfaVerReq)
		if err != nil {
			return nil, err
		}

		var mfaRes authn.FlowVerifyMFAStartResult
		json.Unmarshal([]byte(pangea.Stringify(mfaVerResp.Result)), &mfaRes)

		fmt.Println(mfaRes)

		if mfaRes.Error != nil {
			return nil, errors.New(*mfaRes.Error)
		}

	}

	u := types.User{
		Email:  emailId,
		FlowId: flowId,
	}

	return &u, nil

}

func CompleteLogin(flowId, code string) (*authn.FlowCompleteResult, error) {

	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()
	fmt.Println("Here 1")
	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.FlowVerifyMFACompleteRequest{
		FlowID: flowId,
		Code:   pangea.String(code),
	}

	mfaVerResp, err := authncli.Flow.Verify.MFA.Complete(ctx, input)
	if err != nil {
		return nil, err
	}

	var x authn.FlowVerifyMFACompleteResult
	json.Unmarshal([]byte(pangea.Stringify(mfaVerResp.Result)), &x)

	if x.Error != nil {
		return nil, errors.New(*x.Error)
	}

	fmt.Println(x)

	var fcr authn.FlowCompleteResult
	if x.NextStep == types.LOGIN_COMPLETE_PATH {
		flowCompleteReq := authn.FlowCompleteRequest{
			FlowID: flowId,
		}

		flowCompleteRes, err := authncli.Flow.Complete(ctx, flowCompleteReq)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(pangea.Stringify(flowCompleteRes.Result)), &fcr)

	}

	return &fcr, nil

}

func LogoutFromSession(token string) (*authn.ClientSessionLogoutResult, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()
	fmt.Println("Here 1")
	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.ClientSessionLogoutRequest{
		Token: token,
	}

	resp, err := authncli.Client.Session.Logout(ctx, input)
	if err != nil {
		return nil, err
	}

	var res *authn.ClientSessionLogoutResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &res)

	return res, nil
}
