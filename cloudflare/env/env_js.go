//go:build js && wasm

package env

import (
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
)

// Getenv gets a value of an environment variable.
//   - https://developers.cloudflare.com/workers/platform/environment-variables/
//   - This function panics when a runtime context is not found.
func Getenv(name string) string {
	return cfruntimecontext.MustGetRuntimeContextEnv().Get(name).String()
}

// GetBinding gets a value of an environment binding.
//   - https://developers.cloudflare.com/workers/platform/bindings/about-service-bindings/
//   - This function panics when a runtime context is not found.
func GetBinding(name string) js.Value {
	return cfruntimecontext.MustGetRuntimeContextEnv().Get(name)
}
