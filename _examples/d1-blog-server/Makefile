.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen
	tinygo build -o ./build/app.wasm -target wasm ./main.go

.PHONY: deploy
deploy:
	wrangler deploy

.PHONY: init-db
init-db:
	wrangler d1 execute d1-blog-server --file=./schema.sql

.PHONY: init-db-preview
init-db-preview:
	wrangler d1 execute d1-blog-server-preview --file=./schema.sql
