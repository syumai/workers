# queues

An example of using Cloudflare Workers that interact with [Cloudflare Queues](https://developers.cloudflare.com/queues/).

## Running

### Requirements

This project requires these tools to be installed globally.

* wrangler
* tinygo

### Supported commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make deploy # deploy worker
```

### Interacting with the local queue

1. Start the dev server.
```sh
make dev
```

2. Send a message to the queue.
```sh
curl -v -X POST http://localhost:8787/ -d '{"message": "Hello, World!"}' -H "Content-Type: application/json"
```

3. Observe the response and server logs

4. You can pass `text/plain` content type to write queue message as the string or omit the `Content-Type` header to write queue message as 
byte array.

