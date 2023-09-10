package pangeaauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/sap200/TinyFundConnect/db"
	"github.com/sap200/TinyFundConnect/pangeachecks"
	"github.com/sap200/TinyFundConnect/secret"
	"github.com/sap200/TinyFundConnect/types"
)

func InitiateSignupFlow(emailId, password, firstName, lastName, ip string) (*types.User, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()
	fmt.Println("Here 1")
	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	fmt.Println("here2")
	fts := []authn.FlowType{authn.FTsignup}
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

	// signup password
	signUpRequest := authn.FlowSignupPasswordRequest{
		FlowID:    flowId,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}

	signUpResp, err := authncli.Flow.Signup.Password(ctx, signUpRequest)
	if err != nil {
		return nil, err
	}
	var fspr authn.FlowSignupPasswordResult
	json.Unmarshal([]byte(pangea.Stringify(signUpResp.Result)), &fspr)

	if fspr.CommonFlowResult.NextStep != pangeachecks.VERIFY_EMAIL_PATH {
		return nil, errors.New("Next Path is not verify/email. Invalid flow detected " + pangea.Stringify(signUpResp.Result))
	}

	// get user from pangea and save userId as well
	getUserRequest := authn.UserProfileGetRequest{
		Email: pangea.String(emailId),
	}

	getUserResp, err := authncli.User.Profile.Get(ctx, getUserRequest)
	if err != nil {
		return nil, err
	}

	var ur authn.UserProfileGetResult
	json.Unmarshal([]byte(pangea.Stringify(getUserResp.Result)), &ur)
	fmt.Println(ur)

	// save flow id to db
	var u types.User

	u = types.User{
		Email:         emailId,
		CustomerId:    ur.ID,
		EmailVerified: ur.Verified,
		FirstName:     firstName,
		LastName:      lastName,
		Ip:            ip,
		FlowId:        flowId,
	}
	m, err := u.EncodeToMap()
	if err != nil {
		return nil, errors.New("Error while unmarshalling user struct to map")
	}
	repo := db.New()
	repo.Save(secret.USER_COLLECTIONS, emailId, m)
	fmt.Println("here")

	return &u, nil
}

func ResendEmailVerification(emailId, flowId string) error {
	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	emailVerificationRequest := authn.FlowVerifyEmailRequest{
		FlowID: flowId,
	}

	emaiVerificationResp, err := authncli.Flow.Verify.Email(ctx, emailVerificationRequest)
	if err != nil {
		return err
	}
	var fvr authn.FlowVerifyEmailResult
	json.Unmarshal([]byte(pangea.Stringify(emaiVerificationResp.Result)), &fvr)

	fmt.Println(fvr)

	return nil
}

func GetEmailVerificationStatus(emailId string) (bool, error) {

	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.UserProfileGetRequest{
		Email: pangea.String(emailId),
	}

	resp, err := authncli.User.Profile.Get(ctx, input)
	if err != nil {
		return false, err
	}

	fmt.Println(pangea.Stringify(resp.Result))

	var ur authn.UserProfileGetResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &ur)

	return ur.Verified, nil
}

// enroll MFA
func EnrollMFA(emailId, flowId string) error {
	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.FlowEnrollMFAStartRequest{
		FlowID:      flowId,
		MFAProvider: authn.MFAPEmailOTP,
	}

	resp, err := authncli.Flow.Enroll.MFA.Start(ctx, input)
	if err != nil {
		return err
	}

	var x authn.FlowEnrollMFAStartResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &x)
	fmt.Println(x)
	if x.CommonFlowResult.Error != nil {
		return errors.New(*x.CommonFlowResult.Error)
	}

	return nil
}

// complete mfa enrollment
func CompleteEnrollMFAComplete(emailId, flowId, otp string) error {
	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.FlowEnrollMFACompleteRequest{
		FlowID: flowId,
		Code:   otp,
	}

	resp, err := authncli.Flow.Enroll.MFA.Complete(ctx, input)
	if err != nil {
		return err
	}

	var x authn.FlowEnrollMFACompleteResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &x)
	fmt.Println(x)
	if x.CommonFlowResult.Error != nil {
		return errors.New(*x.CommonFlowResult.Error)
	}

	return nil
}

// complete signup or signin Flow
func CompleteSignupOrLoginFlow(flowId string) (*authn.FlowCompleteResult, error) {

	// send verify email link
	ctx, cancelFn := context.WithTimeout(context.Background(), secret.CANCEL_TIMEOUT_PANGEA*time.Second)
	defer cancelFn()

	// get pangea token and create a client
	authncli := authn.New(&pangea.Config{
		Token:  secret.PANGEA_AUTHN_TOKEN,
		Domain: secret.PANGEA_DOMAIN,
	})

	input := authn.FlowCompleteRequest{
		FlowID: flowId,
	}

	resp, err := authncli.Flow.Complete(ctx, input)
	if err != nil {
		return nil, err
	}

	var x authn.FlowCompleteResult
	json.Unmarshal([]byte(pangea.Stringify(resp.Result)), &x)

	return &x, nil

}
