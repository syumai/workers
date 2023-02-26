.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	mkdir -p dist
	tinygo build -o ./dist/app.wasm -target wasm ./main.go

.PHONY: publish
publish:
	wrangler publish

.PHONY: generate
generate:
	go generate ./...

.PHONY: init-db
init-db:
	wrangler d1 execute d1-blog-server --file=./schema.sql

.PHONY: init-db-preview
init-db-preview:
	wrangler d1 execute d1-blog-server-preview --file=./schema.sql
