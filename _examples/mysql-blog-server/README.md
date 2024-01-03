# mysql-blog-server

* A simple Blog server implemented in Go.
* This example is using MySQL.

# WIP

### Create blog post

```
$ curl -X POST 'http://localhost:8787/articles' \
-H 'Content-Type: application/json' \
-d '{
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
$ curl 'http://localhost:8787/articles'
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
* go

### Setup MySQL DB

* This project requires MySQL DB.
  - Connection setting: `.dev.vars.example` (please rename to `.dev.vars`.)
  - Initial migration SQL: `schema.sql`
* If you want to deploy this app to production, please set `MYSQL_DSN` to your Worker secrets.
  - Run: `npx wrangler secret put MYSQL_DSN`.

### Commands

```
make dev    # run dev server
make build  # build Go Wasm binary
make deploy # deploy worker
```
