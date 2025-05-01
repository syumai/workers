# durable object counter

This app is an exmaple of using a stub to access a direct object. The example
is based on the [cloudflare/durable-object-template](https://github.com/cloudflare/durable-objects-template)
repository.

_The durable object is written in js; only the stub is called from go!_

## Demo

After `make deploy` the trigger is `http://durable-object-counter.YOUR-DOMAIN.workers.dev`

* https://durable-object-counter.YOUR-DOMAIN.workers.dev/
* https://durable-object-counter.YOUR-DOMAIN.workers.dev/increment
* https://durable-object-counter.YOUR-DOMAIN.workers.dev/decrement

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

## Author

akarasz

## License

MIT
