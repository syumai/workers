# wait-until

* Executes a task in worker even after the server has returned the response.
* This example executes 5-second task after responding and outputs the console log.

## Document

* https://developers.cloudflare.com/workers/runtime-apis/fetch-event/#waituntil

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
