package workers

// Getenv gets an environment variable set to global object.
// * https://developers.cloudflare.com/workers/platform/environment-variables/
// * Technically, this function is just an alias for js.Global().Get(env_name).String().
func Getenv(name string) string {
	return global.Get(name).String()
}
