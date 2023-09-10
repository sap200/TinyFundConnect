package pangeaauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
)

func VerifyEmailForPasswordReset(emailId string) (string, error) {
	// start flow
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
		return "", err
	}
	var fr authn.FlowStartResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &fr)
	flowId := fr.FlowID
	if flowId == "" {
		return "", errors.New("Invalid Flow Id: " + pangea.Stringify(resp.Result))
	}
	fmt.Println("FLow Id here: ", flowId)

	// send password reset verify with reset = true
	var shouldReset bool
	shouldReset = true
	passwordResetRequest := authn.FlowVerifyPasswordRequest{
		FlowID: flowId,
		Reset:  &shouldReset,
	}

	passVerRes, err := authncli.Flow.Verify.Password(ctx, passwordResetRequest)
	if err != nil {
		return "", err
	}

	var pr authn.FlowVerifyPasswordResult
	json.Unmarshal([]byte(pangea.Stringify(passVerRes.Result)), &pr)

	if pr.Error != nil {
		return "", errors.New(*pr.Error)
	}

	fmt.Println(pr)

	if pr.NextStep == types.PANGEA_MAIL_VERIFY_PASSWORD_RESET_FINAL_PATH {
		return flowId, nil
	}

	return "", errors.New("Unforeseen Error")
}

func EnterPasswordAndStartMFAVerification(flowId, password string, cancel bool) (string, error) {
	// start flow
	// start the flow
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()
	fmt.Println("Here 1")
	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	if cancel {

		isCancelled := cancelPassLink(flowId)
		if isCancelled {
			return types.PASSWORD_RESET_CANCELLATION_MESSAGE, nil
		}

		return "", errors.New("Unable to cancel password reset link! Please try to verify the email or it may be possible that this is an invalid flow now !")
	}

	input := authn.FlowResetPasswordRequest{
		FlowID:   flowId,
		Password: password,
	}

	resp, err := authncli.Flow.Reset.Password(ctx, input)
	if err != nil {
		return "", nil
	}

	var pr authn.FlowResetPasswordResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &pr)
	if pr.Error != nil {
		return "", errors.New(*pr.Error)
	}

	if pr.NextStep == types.MFA_VERIFICATION_PATH {
		mfaVerReq := authn.FlowVerifyMFAStartRequest{
			FlowID:      flowId,
			MFAProvider: authn.MFAPEmailOTP,
		}

		mfaVerResp, err := authncli.Flow.Verify.MFA.Start(ctx, mfaVerReq)
		if err != nil {
			return "", err
		}

		var mfaRes authn.FlowVerifyMFAStartResult
		json.Unmarshal([]byte(pangea.Stringify(mfaVerResp.Result)), &mfaRes)

		fmt.Println(mfaRes)

		if mfaRes.Error != nil {
			return "", errors.New(*mfaRes.Error)
		}

		return flowId, nil
	}

	return "", errors.New(*pr.Error)
}

func ConfirmPasswordChangeWithMFAVerification(flowId, code string) (bool, error) {
	// start flow
	// start the flow
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
		return false, err
	}

	var x authn.FlowVerifyMFACompleteResult
	json.Unmarshal([]byte(pangea.Stringify(mfaVerResp.Result)), &x)

	if x.Error != nil {
		return false, errors.New(*x.Error)
	}

	fmt.Println(x)

	var fcr authn.FlowCompleteResult
	if x.NextStep == types.LOGIN_COMPLETE_PATH {
		flowCompleteReq := authn.FlowCompleteRequest{
			FlowID: flowId,
		}

		flowCompleteRes, err := authncli.Flow.Complete(ctx, flowCompleteReq)
		if err != nil {
			return false, err
		}

		json.Unmarshal([]byte(pangea.Stringify(flowCompleteRes.Result)), &fcr)

	}

	fmt.Println("Active Token:: " + fcr.ActiveToken.Token)

	if fcr.ActiveToken.Token != "" {
		return true, nil
	}

	return false, errors.New("Unknown error occured")

}

func ShouldWeCircleInVerifyEmailPage(flowId string) (bool, error) {

	url := "https://authn.aws.eu.pangea.cloud/v1/flow/get"
	method := "POST"

	payload := strings.NewReader("{\"flow_id\": \"" + flowId + "\"}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Authorization", "Bearer "+secret.PANGEA_AUTHN_TOKEN)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "AWSALB=lxwlOcqqVI4zR0XTZQJyOMVuyIC+uNdMxOLw6ReiDPIvTVQp9O8cS0a33PQey87/W9qCt3WMSZiIIjJToZ3uJv1IT/aF2adKD7XMh03K0l72LuzR0II47hfLOjWc; AWSALBCORS=lxwlOcqqVI4zR0XTZQJyOMVuyIC+uNdMxOLw6ReiDPIvTVQp9O8cS0a33PQey87/W9qCt3WMSZiIIjJToZ3uJv1IT/aF2adKD7XMh03K0l72LuzR0II47hfLOjWc")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	var response types.FlowResponseModel
	json.Unmarshal(body, &response)

	if response.Result.NextStep == types.PANGEA_MAIL_VERIFY_PASSWORD_RESET_FINAL_PATH {
		return true, nil
	} else if response.Result.NextStep == types.PANGEA_MOVE_FROM_EMAIL_VERIFICATION_PATH_PASSWORD_RESET {
		return false, nil
	}

	fmt.Println(string(body))
	fmt.Println(response)

	return false, errors.New("have received unexpected next path " + response.Result.NextStep)
}

func cancelPassLink(flowId string) bool {
	url := "https://authn.aws.eu.pangea.cloud/v1/flow/reset/password"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("{\"cancel\":true,\"flow_id\": \"%s\"}", flowId))
	fmt.Println(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Add("Authorization", "Bearer pts_ozksat37mioqrcj7djlvn3pxnn6tfuhz")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "AWSALB=M0UpoLS8zmfZc9rP1ytNRngmiFNYSR6Zb/LbiMfBYnMoFDs9axsSVXYw0OKssJOuDz4xy/4JtGn4cZmzuoGJkfGfEC3TZkc3xDdHqrhQCCXB8ii26DSgVujFKWpM; AWSALBCORS=M0UpoLS8zmfZc9rP1ytNRngmiFNYSR6Zb/LbiMfBYnMoFDs9axsSVXYw0OKssJOuDz4xy/4JtGn4cZmzuoGJkfGfEC3TZkc3xDdHqrhQCCXB8ii26DSgVujFKWpM")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println(string(body))

	return strings.Contains(string(body), "Success")
}
