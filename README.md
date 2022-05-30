# workers [![Go Reference](https://pkg.go.dev/badge/github.com/syumai/workers.svg)](https://pkg.go.dev/github.com/syumai/workers)

* `workers` is a package to run an HTTP server written in Go on [Cloudflare Workers](https://workers.cloudflare.com/).
* This package can easily serve *http.Handler* on Cloudflare Workers.
* Caution: This is an experimental project.

## Features

* [x] serve http.Handler
* [ ] R2 - Partially supported
  - [x] Head
  - [x] Get
  - [x] Put (load all bytes to memory)
  - [ ] Put (stream)
  - [x] Delete
  - [x] List
  - [ ] Options for R2 methods
* [ ] environment variables (WIP)
* [ ] KV (WIP)

## Installation

```
go get github.com/syumai/workers
```

## Usage

implement your http.Handler and give it to `workers.Serve()`.

```go
func main() {
	var handler http.HandlerFunc = func (w http.ResponseWriter, req *http.Request) { ... }
	workers.Serve(handler)
}
```

or just call `http.Handle` and `http.HandleFunc`, then invoke `workers.Serve()` with nil.

```go
func main() {
	http.HandleFunc("/hello", func (w http.ResponseWriter, req *http.Request) { ... })
	workers.Serve(nil) // if nil is given, http.DefaultServeMux is used.
}
```

For concrete examples, see `examples` directory.
Currently, all examples use tinygo instead of Go due to binary size issues.

A template repository is also available.
* https://github.com/syumai/worker-template-tinygo

## FAQ

### How do I deploy a worker implemented in this package?

To deploy a Worker, the following steps are required.

* Create a worker project using [wrangler](https://developers.cloudflare.com/workers/wrangler/).
* Build a Wasm binary.
* Upload a Wasm binary with a JavaScript code to load and instantiate Wasm (for entry point).

The [worker-template-tinygo](https://github.com/syumai/worker-template-tinygo) repository contains all the required files, so I recommend using this template.

## License

MIT

## Author

syumai
