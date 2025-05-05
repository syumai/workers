.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm ./...

.PHONY: create-kv-namespace
create-kv-namespace:
	wrangler kv namespace create COUNTER

.PHONY: deploy
deploy:
	wrangler deploy
