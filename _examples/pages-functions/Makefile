.PHONY: dev
dev:
	wrangler pages dev ./pages

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen -mode=go
	GOOS=js GOARCH=wasm go build -o ./build/app.wasm ./...

.PHONY: deploy
deploy:
	wrangler pages deploy ./pages
