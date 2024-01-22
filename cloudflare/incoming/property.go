package incoming

import (
	"context"
	"errors"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
	"github.com/syumai/workers/internal/runtimecontext"
)

type BotManagementJsDetection struct {
	Passed bool
}

func NewBotManagementJsDetection(cf js.Value) *BotManagementJsDetection {
	if cf.IsUndefined() {
		return nil
	}
	return &BotManagementJsDetection{
		Passed: cf.Get("passed").Bool(),
	}
}

type BotManagement struct {
	CorporateProxy bool
	VerifiedBot    bool
	JsDetection    *BotManagementJsDetection
	StaticResource bool
	Score          int
}

func NewBotManagement(cf js.Value) *BotManagement {
	if cf.IsUndefined() {
		return nil
	}
	return &BotManagement{
		CorporateProxy: cf.Get("corporateProxy").Bool(),
		VerifiedBot:    cf.Get("verifiedBot").Bool(),
		JsDetection:    NewBotManagementJsDetection(cf.Get("jsDetection")),
		StaticResource: cf.Get("staticResource").Bool(),
		Score:          cf.Get("score").Int(),
	}
}

type TLSClientAuth struct {
	CertIssuerDNLegacy    string
	CertIssuerSKI         string
	CertSubjectDNRFC2253  string
	CertSubjectDNLegacy   string
	CertFingerprintSHA256 string
	CertNotBefore         string
	CertSKI               string
	CertSerial            string
	CertIssuerDN          string
	CertVerified          string
	CertNotAfter          string
	CertSubjectDN         string
	CertPresented         string
	CertRevoked           string
	CertIssuerSerial      string
	CertIssuerDNRFC2253   string
	CertFingerprintSHA1   string
}

func NewTLSClientAuth(cf js.Value) *TLSClientAuth {
	if cf.IsUndefined() {
		return nil
	}
	return &TLSClientAuth{
		CertIssuerDNLegacy:    jsutil.MaybeString(cf.Get("certIssuerDNLegacy")),
		CertIssuerSKI:         jsutil.MaybeString(cf.Get("certIssuerSKI")),
		CertSubjectDNRFC2253:  jsutil.MaybeString(cf.Get("certSubjectDNRFC2253")),
		CertSubjectDNLegacy:   jsutil.MaybeString(cf.Get("certSubjectDNLegacy")),
		CertFingerprintSHA256: jsutil.MaybeString(cf.Get("certFingerprintSHA256")),
		CertNotBefore:         jsutil.MaybeString(cf.Get("certNotBefore")),
		CertSKI:               jsutil.MaybeString(cf.Get("certSKI")),
		CertSerial:            jsutil.MaybeString(cf.Get("certSerial")),
		CertIssuerDN:          jsutil.MaybeString(cf.Get("certIssuerDN")),
		CertVerified:          jsutil.MaybeString(cf.Get("certVerified")),
		CertNotAfter:          jsutil.MaybeString(cf.Get("certNotAfter")),
		CertSubjectDN:         jsutil.MaybeString(cf.Get("certSubjectDN")),
		CertPresented:         jsutil.MaybeString(cf.Get("certPresented")),
		CertRevoked:           jsutil.MaybeString(cf.Get("certRevoked")),
		CertIssuerSerial:      jsutil.MaybeString(cf.Get("certIssuerSerial")),
		CertIssuerDNRFC2253:   jsutil.MaybeString(cf.Get("certIssuerDNRFC2253")),
		CertFingerprintSHA1:   jsutil.MaybeString(cf.Get("certFingerprintSHA1")),
	}
}

type TLSExportedAuthenticator struct {
	ClientFinished  string
	ClientHandshake string
	ServerHandshake string
	ServerFinished  string
}

func NewTLSExportedAuthenticator(cf js.Value) *TLSExportedAuthenticator {
	if cf.IsUndefined() {
		return nil
	}
	return &TLSExportedAuthenticator{
		ClientFinished:  jsutil.MaybeString(cf.Get("clientFinished")),
		ClientHandshake: jsutil.MaybeString(cf.Get("clientHandshake")),
		ServerHandshake: jsutil.MaybeString(cf.Get("serverHandshake")),
		ServerFinished:  jsutil.MaybeString(cf.Get("serverFinished")),
	}
}

type Properties struct {
	Longitude                string
	Latitude                 string
	TLSCipher                string
	Continent                string
	Asn                      int
	ClientAcceptEncoding     string
	Country                  string
	TLSClientAuth            *TLSClientAuth
	TLSExportedAuthenticator *TLSExportedAuthenticator
	TLSVersion               string
	Colo                     string
	Timezone                 string
	City                     string
	VerifiedBotCategory      string
	// EdgeRequestKeepAliveStatus int
	RequestPriority string
	HttpProtocol    string
	Region          string
	RegionCode      string
	AsOrganization  string
	PostalCode      string
	BotManagement   *BotManagement
}

func NewProperties(ctx context.Context) (*Properties, error) {
	obj := runtimecontext.MustExtractTriggerObj(ctx)
	cf := obj.Get("cf")
	if cf.IsUndefined() {
		return nil, errors.New("runtime is not cloudflare")
	}

	return &Properties{
		Longitude:                jsutil.MaybeString(cf.Get("longitude")),
		Latitude:                 jsutil.MaybeString(cf.Get("latitude")),
		TLSCipher:                jsutil.MaybeString(cf.Get("tlsCipher")),
		Continent:                jsutil.MaybeString(cf.Get("continent")),
		Asn:                      cf.Get("asn").Int(),
		ClientAcceptEncoding:     jsutil.MaybeString(cf.Get("clientAcceptEncoding")),
		Country:                  jsutil.MaybeString(cf.Get("country")),
		TLSClientAuth:            NewTLSClientAuth(cf.Get("tlsClientAuth")),
		TLSExportedAuthenticator: NewTLSExportedAuthenticator(cf.Get("tlsExportedAuthenticator")),
		TLSVersion:               cf.Get("tlsVersion").String(),
		Colo:                     cf.Get("colo").String(),
		Timezone:                 cf.Get("timezone").String(),
		City:                     jsutil.MaybeString(cf.Get("city")),
		VerifiedBotCategory:      jsutil.MaybeString(cf.Get("verifiedBotCategory")),
		RequestPriority:          jsutil.MaybeString(cf.Get("requestPriority")),
		HttpProtocol:             cf.Get("httpProtocol").String(),
		Region:                   jsutil.MaybeString(cf.Get("region")),
		RegionCode:               jsutil.MaybeString(cf.Get("regionCode")),
		AsOrganization:           cf.Get("asOrganization").String(),
		PostalCode:               jsutil.MaybeString(cf.Get("postalCode")),
		BotManagement:            NewBotManagement(cf.Get("botManagement")),
	}, nil
}
