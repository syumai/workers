package workers

// Getenv is an alias for js.Global().Get(env_name).
func Getenv(name string) string {
	return global.Get(name).String()
}
