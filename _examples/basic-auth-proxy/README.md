# basic-auth-proxy

* This is an example of an HTTP proxy server with Basic-Auth .
* This proxy server adds Basic-Auth to `https://syum.ai` .

## Demo

* https://basic-auth-proxy.syumai.workers.dev/
* Try:
  - userName: `user`
  - password: `password`

## Development

### Requirements

This project requires these tools to be installed globally.

* wrangler
* tinygo

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy # deploy worker
```
