# simple-json-server

* A simple HTTP JSON server implemented in Go and compiled with tinygo.

## Example

* https://simple-json-server.syumai.workers.dev

### Request

```
curl --location --request POST 'https://simple-json-server.syumai.workers.dev/hello' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "syumai"
}'
```

### Response

```json
{
    "message": "Hello, syumai!"
}
```

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* tinygo
* [easyjson](https://github.com/mailru/easyjson)
  - `go install github.com/mailru/easyjson/...@latest`

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make publish # publish worker
```

## Author

syumai

## License

MIT
