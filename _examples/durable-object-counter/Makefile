.PHONY: dev
dev:
	wrangler dev

.PHONY: build
build:
	mkdir -p build
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy:
	wrangler deploy
