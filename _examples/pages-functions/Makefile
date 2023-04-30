.PHONY: dev
dev:
	wrangler pages dev ./pages

.PHONY: build
build:
	mkdir -p build
	tinygo build -o ./build/app.wasm -target wasm ./...

.PHONY: publish
publish:
	wrangler pages publish ./pages
