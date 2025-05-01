.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm .

.PHONY: deploy
deploy:
	wrangler deploy

.PHONY: create-db
create-db:
	wrangler d1 create d1-blog-server

.PHONY: init-db
init-db:
	wrangler d1 execute d1-blog-server --remote --file=./schema.sql

.PHONY: init-db-local
init-db-local:
	wrangler d1 execute d1-blog-server --file=./schema.sql --local
