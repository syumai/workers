# hello

* This app just returns a message `Hello, world!`.
* If url param like `?name=syumai`, then a message `Hello, syumai!` will be returned.

## Demo

* https://hello.syumai.workers.dev/
* https://hello.syumai.workers.dev/?name=syumai

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

## Author

syumai

## License

MIT
