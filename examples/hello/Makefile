.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run ../../cmd/workers-gen-go
	tinygo build -o ./build/app.wasm -target wasm ./...

.PHONY: publish
publish:
	wrangler publish
