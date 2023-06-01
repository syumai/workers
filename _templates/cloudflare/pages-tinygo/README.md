# pages-tinygo

- A template for starting a Cloudflare Pages Functions project with tinygo.
- This template uses the [`workers`](https://github.com/syumai/workers) package to run.

## Usage

- `main.go` includes a [chi](https://github.com/go-chi/chi) HTTP router implementation with three different routes. Feel free to edit this code and implement your own HTTP router.

## Requirements

- Node.js
- [wrangler](https://developers.cloudflare.com/workers/wrangler/)
  - just run `npm install -g wrangler`
- tinygo

## Getting Started

```console
wrangler generate my-app syumai/workers/_templates/cloudflare/pages-tinygo
cd my-app
go mod init
go mod tidy
make dev # start running dev server
curl http://localhost:8787/api/hello # outputs "Hello, Pages Functions!"
```

## Development

### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy # deploy worker
```

### Testing dev server

- You can send HTTP requests using tools like curl.

```
$ curl http://localhost:8787/api/hello
Hello, Pages Functions!
```

```
$ curl http://localhost:8787/api/hello?name=Example
Hello, Example!
```

```
$ curl http://localhost:8787/api/hello2
Hello, Hello world!
```

```
$ curl http://localhost:8787/api/hello3
Hello, Hello, Hello world!
```

