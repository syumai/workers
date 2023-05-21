.PHONY: dev
dev:
	wrangler pages dev ./pages

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen
	tinygo build -o ./build/app.wasm -target wasm ./...

.PHONY: deploy
deploy:
	wrangler pages deploy ./pages
