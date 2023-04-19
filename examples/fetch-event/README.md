# fetch-event

## Document

* https://developers.cloudflare.com/workers/runtime-apis/fetch-event/

## Setup

```toml
routes = [
    { pattern = "example.com/*", zone_name = "example.com" }
]
```

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* tinygo

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make publish # publish worker
```
