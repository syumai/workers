.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm ./...

.PHONY: create-queue
create-queue:
	wrangler queues create my-queue

.PHONY: deploy
deploy:
	wrangler deploy
