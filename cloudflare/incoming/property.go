package incoming

import (
	"context"
	"syscall/js"

	"github.com/syumai/workers/internal/cfcontext"
	"github.com/syumai/workers/internal/jsutil"
)

type BotManagementJsDetection struct {
	Passed bool
}

func NewBotManagementJsDetection(obj js.Value) *BotManagementJsDetection {
	if obj.IsUndefined() {
		return nil
	}
	return &BotManagementJsDetection{
		Passed: obj.Get("passed").Bool(),
	}
}

type BotManagement struct {
	CorporateProxy bool
	VerifiedBot    bool
	JsDetection    *BotManagementJsDetection
	StaticResource bool
	Score          int
}

func NewBotManagement(obj js.Value) *BotManagement {
	if obj.IsUndefined() {
		return nil
	}
	return &BotManagement{
		CorporateProxy: obj.Get("corporateProxy").Bool(),
		VerifiedBot:    obj.Get("verifiedBot").Bool(),
		JsDetection:    NewBotManagementJsDetection(obj.Get("jsDetection")),
		StaticResource: obj.Get("staticResource").Bool(),
		Score:          obj.Get("score").Int(),
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

func NewTLSClientAuth(obj js.Value) *TLSClientAuth {
	if obj.IsUndefined() {
		return nil
	}
	return &TLSClientAuth{
		CertIssuerDNLegacy:    jsutil.MaybeString(obj.Get("certIssuerDNLegacy")),
		CertIssuerSKI:         jsutil.MaybeString(obj.Get("certIssuerSKI")),
		CertSubjectDNRFC2253:  jsutil.MaybeString(obj.Get("certSubjectDNRFC2253")),
		CertSubjectDNLegacy:   jsutil.MaybeString(obj.Get("certSubjectDNLegacy")),
		CertFingerprintSHA256: jsutil.MaybeString(obj.Get("certFingerprintSHA256")),
		CertNotBefore:         jsutil.MaybeString(obj.Get("certNotBefore")),
		CertSKI:               jsutil.MaybeString(obj.Get("certSKI")),
		CertSerial:            jsutil.MaybeString(obj.Get("certSerial")),
		CertIssuerDN:          jsutil.MaybeString(obj.Get("certIssuerDN")),
		CertVerified:          jsutil.MaybeString(obj.Get("certVerified")),
		CertNotAfter:          jsutil.MaybeString(obj.Get("certNotAfter")),
		CertSubjectDN:         jsutil.MaybeString(obj.Get("certSubjectDN")),
		CertPresented:         jsutil.MaybeString(obj.Get("certPresented")),
		CertRevoked:           jsutil.MaybeString(obj.Get("certRevoked")),
		CertIssuerSerial:      jsutil.MaybeString(obj.Get("certIssuerSerial")),
		CertIssuerDNRFC2253:   jsutil.MaybeString(obj.Get("certIssuerDNRFC2253")),
		CertFingerprintSHA1:   jsutil.MaybeString(obj.Get("certFingerprintSHA1")),
	}
}

type TLSExportedAuthenticator struct {
	ClientFinished  string
	ClientHandshake string
	ServerHandshake string
	ServerFinished  string
}

func NewTLSExportedAuthenticator(obj js.Value) *TLSExportedAuthenticator {
	if obj.IsUndefined() {
		return nil
	}
	return &TLSExportedAuthenticator{
		ClientFinished:  jsutil.MaybeString(obj.Get("clientFinished")),
		ClientHandshake: jsutil.MaybeString(obj.Get("clientHandshake")),
		ServerHandshake: jsutil.MaybeString(obj.Get("serverHandshake")),
		ServerFinished:  jsutil.MaybeString(obj.Get("serverFinished")),
	}
}

type Properties struct {
	Longitude                string
	Latitude                 string
	TlsCipher                string
	Continent                string
	Asn                      int
	ClientAcceptEncoding     string
	Country                  string
	TLSClientAuth            *TLSClientAuth
	TLSExportedAuthenticator *TLSExportedAuthenticator
	TlsVersion               string
	Colo                     string
	Timezone                 string
	City                     string
	VerifiedBotCategory      string
	// EdgeRequestKeepAliveStatus int
	RequestPriority string `json:""`
	HttpProtocol    string `json:""`
	Region          string `json:"region"`
	RegionCode      string `json:"regionCode"`
	AsOrganization  string `json:"asOrganization"`
	PostalCode      string `json:"postalCode"`
	BotManagement   *BotManagement
}

func NewProperties(ctx context.Context) *Properties {
	obj := cfcontext.MustExtractIncomingProperty(ctx)
	return &Properties{
		Longitude:                jsutil.MaybeString(obj.Get("longitude")),
		Latitude:                 jsutil.MaybeString(obj.Get("latitude")),
		TlsCipher:                jsutil.MaybeString(obj.Get("tlsCipher")),
		Continent:                jsutil.MaybeString(obj.Get("continent")),
		Asn:                      obj.Get("asn").Int(),
		ClientAcceptEncoding:     jsutil.MaybeString(obj.Get("clientAcceptEncoding")),
		Country:                  jsutil.MaybeString(obj.Get("country")),
		TLSClientAuth:            NewTLSClientAuth(obj.Get("tlsClientAuth")),
		TLSExportedAuthenticator: NewTLSExportedAuthenticator(obj.Get("tlsExportedAuthenticator")),
		TlsVersion:               obj.Get("tlsVersion").String(),
		Colo:                     obj.Get("colo").String(),
		Timezone:                 obj.Get("timezone").String(),
		City:                     jsutil.MaybeString(obj.Get("city")),
		VerifiedBotCategory:      jsutil.MaybeString(obj.Get("verifiedBotCategory")),
		RequestPriority:          jsutil.MaybeString(obj.Get("requestPriority")),
		HttpProtocol:             obj.Get("httpProtocol").String(),
		Region:                   jsutil.MaybeString(obj.Get("region")),
		RegionCode:               jsutil.MaybeString(obj.Get("regionCode")),
		AsOrganization:           obj.Get("asOrganization").String(),
		PostalCode:               jsutil.MaybeString(obj.Get("postalCode")),
		BotManagement:            NewBotManagement(obj.Get("botManagement")),
	}
}
