package fetch

import (
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

// RedirectMode represents the redirect mode of a fetch() request.
type RedirectMode string

var (
	RedirectModeFollow RedirectMode = "follow"
	RedirectModeError  RedirectMode = "error"
	RedirectModeManual RedirectMode = "manual"
)

func (mode RedirectMode) IsValid() bool {
	return mode == RedirectModeFollow || mode == RedirectModeError || mode == RedirectModeManual
}

func (mode RedirectMode) String() string {
	return string(mode)
}

// RequestInit represents the options passed to a fetch() request.
type RequestInit struct {
	CF       *RequestInitCF
	Redirect RedirectMode
}

// ToJS converts RequestInit to JS object.
func (init *RequestInit) ToJS() js.Value {
	if init == nil {
		return js.Undefined()
	}
	obj := jsutil.NewObject()
	if init.Redirect.IsValid() {
		obj.Set("redirect", init.Redirect.String())
	}
	return obj
}

// RequestInitCF represents the Cloudflare-specific options passed to a fetch() request.
type RequestInitCF struct {
	/* TODO: implement */
}
