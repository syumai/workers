// Deprecated: use github.com/syumai/workers-go/cloudflare/fetch instead.
package fetch

import (
	src "github.com/syumai/workers-go/cloudflare/fetch"
)

type Client = src.Client
type ClientOption = src.ClientOption
type IncomingBotManagement = src.IncomingBotManagement
type IncomingBotManagementJsDetection = src.IncomingBotManagementJsDetection
type IncomingProperties = src.IncomingProperties
type IncomingTLSClientAuth = src.IncomingTLSClientAuth
type IncomingTLSExportedAuthenticator = src.IncomingTLSExportedAuthenticator

var NewClient = src.NewClient
var NewIncomingBotManagement = src.NewIncomingBotManagement
var NewIncomingBotManagementJsDetection = src.NewIncomingBotManagementJsDetection
var NewIncomingProperties = src.NewIncomingProperties
var NewIncomingTLSClientAuth = src.NewIncomingTLSClientAuth
var NewIncomingTLSExportedAuthenticator = src.NewIncomingTLSExportedAuthenticator
var NewRequest = src.NewRequest

type RedirectMode = src.RedirectMode

var RedirectModeError = src.RedirectModeError
var RedirectModeFollow = src.RedirectModeFollow
var RedirectModeManual = src.RedirectModeManual

type Request = src.Request
type RequestInit = src.RequestInit
type RequestInitCF = src.RequestInitCF

var WithBinding = src.WithBinding
