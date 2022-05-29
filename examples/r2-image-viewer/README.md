# r2-image-viewer-tinygo

* An example server which returns image from Cloudflare R2.
* This server is implemented in Go and compiled with tinygo.

## Example

* https://r2-image-viewer-tinygo.syumai.workers.dev/syumai.png

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
