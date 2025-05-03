# r2-image-viewer-tinygo

* An example server which returns image from Cloudflare R2.
* This server is implemented in Go.

## Example

* https://r2-image-viewer-tinygo.syumai.workers.dev/syumai.png

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* Go 1.24.0 or later

### Commands

```
make dev           # run dev server
make build         # build Go Wasm binary
make create-bucket # create r2 bucket
make deploy        # deploy worker
```
