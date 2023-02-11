package cloudflare

import (
	"context"
)

// Getenv gets a value of an environment variable.
//   - https://developers.cloudflare.com/workers/platform/environment-variables/
//   - This function panics when a runtime context is not found.
func Getenv(ctx context.Context, name string) string {
	return getRuntimeContextEnv(ctx).Get(name).String()
}
