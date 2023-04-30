# [FetchEvent](https://developers.cloudflare.com/workers/runtime-apis/fetch-event/)

Normally, workers are designed to return some kind of HTTP Response and exit immediately upon receiving an HTTP request. `FetchEvent` can extend these life cycles.

#### WaitUntil

`WaitUntil` extends the lifetime of the "fetch" event. It accepts an asynchronous task which the Workers runtime will execute without blocking the response. The worker will not be terminated until those tasks are completed.

#### PassThroughOnException

`PassThroughOnException` prevents a runtime error response when the Worker script throws an unhandled exception. Instead, the request forwards to the origin server as if it had not gone through the worker.

## Example

### Usecase

You have decided to implement a log stream API to capture all access logs. You must edit the Headers so that the user's API token is not logged. If an unknown error occurs during this process, the entire service will be down, which must be avoided.
In such cases, declare PassThroughOnException first and use WaitUntil for logging.

### Setup

This example worker is triggered by [Routes](https://developers.cloudflare.com/workers/platform/triggers/routes/). To try this example, add your site to cloudflare and add some records(A and CNAME, etc.) so that you can actually access the website.
If your domain is `sub.example.com`, edit `wrangler.toml` as following:

```toml
routes = [
    { pattern = "sub.example.com/*", zone_name = "example.com" }
]
```

The workers is executed if the URL matches `sub.example.com/*`.

### Development

#### Requirements

This project requires these tools to be installed globally.

* wrangler
* tinygo

#### Commands

```
make dev     # run dev server
make build   # build Go Wasm binary
make publish # publish worker
```