.PHONY: dev
dev:
	wrangler pages dev ./pages

.PHONY: build
build:
	go run github.com/syumai/workers/cmd/workers-assets-gen
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy:
	wrangler pages deploy ./pages
