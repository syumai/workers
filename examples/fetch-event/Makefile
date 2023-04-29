.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen@latest
	#tinygo build -o ./dist/app.wasm -target wasm ./...
	tinygo build -o ./dist/app.wasm -target wasm -no-debug ./...

.PHONY: publish
publish:
	wrangler publish
