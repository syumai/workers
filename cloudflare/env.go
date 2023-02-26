package cloudflare

import (
	"context"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
)

// Getenv gets a value of an environment variable.
//   - https://developers.cloudflare.com/workers/platform/environment-variables/
//   - This function panics when a runtime context is not found.
func Getenv(ctx context.Context, name string) string {
	return cfruntimecontext.GetRuntimeContextEnv(ctx).Get(name).String()
}
