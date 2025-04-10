package main

type Runtime string

const (
	RuntimeCloudflare Runtime = "cloudflare"
)

func (r Runtime) IsValid() bool {
	switch r {
	case RuntimeCloudflare:
		return true
	}
	return false
}

func (r Runtime) AssetFileName() string {
	return string(r) + ".mjs"
}
