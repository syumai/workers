.PHONY: dev
dev:
	wrangler pages dev ./pages

.PHONY: build
build:
	mkdir -p dist
	tinygo build -o ./dist/app.wasm -target wasm ./...

.PHONY: publish
publish:
	wrangler pages publish ./pages
