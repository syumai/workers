.PHONY: dev
dev:
	wrangler dev

.PHONY: build-tinygo
build-tinygo:
	mkdir -p dist
	tinygo build -o ./dist/app.wasm -target wasm ./main.go

.PHONY: build
build:
	mkdir -p dist
	GOOS=js GOARCH=wasm go build -o ./dist/app.wasm .

.PHONY: publish
publish:
	wrangler publish
