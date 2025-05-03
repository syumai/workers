# kv-counter

* This app counts page view using Cloudflare KV.

## Demo

* https://kv-counter.syumai.workers.dev/

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* Go 1.24.0 or later

### Commands

```
make dev                  # run dev server
make build                # build Go Wasm binary
make create-kv-namespace  # creates a kv namespace - add binding to wrangler.toml before deploy
make deploy               # deploy worker
```
