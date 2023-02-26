# d1-blog-server

* A simple Blog server implemented in Go and compiled with tinygo.
* This example is using Cloudflare D1.

# WIP

## Example

* https://d1-blog-server.syumai.workers.dev

### Create blog post

```
$ curl --location --request POST 'https://d1-blog-server.syumai.workers.dev/articles' \
--header 'Content-Type: application/json' \
--data-raw '{
  "title":"example post",
  "body":"body of the example post"
}'
{
  "article": {
    {
      "id": "f9e8119e-881e-4dc5-9307-af4f2dc79891",
      "title": "example post",
      "body": "body of the example post",
      "createdAt": 1677382874
    }
  }
}
```

### List blog posts

```
$ curl 'https://d1-blog-server.syumai.workers.dev/articles'
{
  "articles": [
    {
      "id": "bea6cd80-5a83-45f0-b061-0e13a2ad5fba",
      "title": "example post 2",
      "body": "body of the example post 2",
      "createdAt": 1677383758
    },
    {
      "id": "f9e8119e-881e-4dc5-9307-af4f2dc79891",
      "title": "example post",
      "body": "body of the example post",
      "createdAt": 1677382874
    }
  ]
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

* Before development, 1. create your own D1 database, 2. set database ID to wrangler.toml and run `wrangler d1 migrations apply [DB Name]`.

```
make dev     # run dev server
make build   # build Go Wasm binary
make publish # publish worker
```

## Author

syumai

## License

MIT
