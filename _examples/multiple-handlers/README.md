# multiple-handlers

* This example shows how to use multiple handlers in a single worker.

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* Go 1.24.0 or later

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy  # deploy worker
```

#### Testing cron schedule

* With curl command below, you can test the cron schedule.
  - see: https://developers.cloudflare.com/workers/runtime-apis/handlers/scheduled/#background

```
curl "http://localhost:8787/__scheduled?cron=*+*+*+*+*"
```
