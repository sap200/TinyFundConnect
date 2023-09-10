package pangeachecks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/ip_intel"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/user_intel"

	"github.com/sap200/TinyFundConnect/secret"
)

var checkedEmailList = make(map[string]string)
var checkedIpListForVPNProxy = make(map[string]string)

const (
	EMAIL_IS_BREACHED           = "breached"
	EMAIL_IS_NOT_BREACHED       = "unbreached"
	IP_IS_VPNPROXYMALICIOUS     = "vpn_proxy_malicious"
	IP_IS_NOT_VPNPROXYMALICIOUS = "not_a_vpn_proxy_malicious"
)

func IsUserEmailInBreachedList(email string) (bool, error) {

	state, ok := checkedEmailList[email]
	if ok {
		if state == EMAIL_IS_BREACHED {
			return true, nil
		}

		return false, nil
	}

	token := secret.PANGEA_INTEL_TOKEN

	intelcli := user_intel.New(&pangea.Config{
		Token:  token,
		Domain: secret.PANGEA_DOMAIN,
	})

	ctx := context.Background()
	input := &user_intel.UserBreachedRequest{
		Email:    email,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "spycloud",
	}

	resp, err := intelcli.UserBreached(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	var uintel user_intel.UserBreachedData
	json.Unmarshal([]byte(pangea.Stringify(resp.Result.Data)), &uintel)

	if uintel.FoundInBreach {
		checkedEmailList[email] = EMAIL_IS_BREACHED
	} else {
		checkedEmailList[email] = EMAIL_IS_NOT_BREACHED
	}

	return uintel.FoundInBreach, nil
}

// Check if VPN proxy or Malicous
func IsTheIpVpnOrProxyOrMalicious(ip string) (bool, error) {

	state, ok := checkedIpListForVPNProxy[ip]
	if ok {
		if state == IP_IS_VPNPROXYMALICIOUS {
			return true, nil
		}

		return false, nil
	}

	token := secret.PANGEA_INTEL_TOKEN

	intelcli := ip_intel.New(&pangea.Config{
		Token:  token,
		Domain: secret.PANGEA_DOMAIN,
	})

	//check proxy
	ctx := context.Background()
	input := &ip_intel.IpProxyRequest{
		Ip:       ip,
		Raw:      pangea.Bool(false),
		Verbose:  pangea.Bool(false),
		Provider: "digitalelement",
	}

	resp, err := intelcli.IsProxy(ctx, input)
	if err != nil {
		return false, err
	}

	isIpProxy := resp.Result.Data.IsProxy

	// check vpn

	inp1 := &ip_intel.IpVPNRequest{
		Ip:       ip,
		Raw:      pangea.Bool(false),
		Verbose:  pangea.Bool(false),
		Provider: "digitalelement",
	}

	resp1, err := intelcli.IsVPN(ctx, inp1)
	if err != nil {
		return false, err
	}

	isIpVpn := resp1.Result.Data.IsVPN

	inp := &ip_intel.IpReputationRequest{
		Ip:       ip,
		Raw:      pangea.Bool(true),
		Verbose:  pangea.Bool(true),
		Provider: "crowdstrike",
	}

	resp2, err := intelcli.Reputation(ctx, inp)
	if err != nil {
		return false, err
	}

	var reputation ip_intel.ReputationData
	json.Unmarshal([]byte(pangea.Stringify(resp2.Result.Data)), &reputation)
	isMalicious := reputation.Verdict == IP_MALICIOUS_CONSTANT

	fmt.Println("IsProxy: ", isIpProxy, " IsVpn: ", isIpVpn, " IsMalicious: ", isMalicious, " Reputation Data: ", reputation)

	if isIpProxy || isIpVpn || isMalicious {
		checkedIpListForVPNProxy[ip] = IP_IS_VPNPROXYMALICIOUS
	} else {
		checkedIpListForVPNProxy[ip] = IP_IS_NOT_VPNPROXYMALICIOUS
	}

	return isIpProxy || isIpVpn || isMalicious, nil

}
