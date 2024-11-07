.PHONY: dev
dev:
	npx wrangler dev --port 8787

.PHONY: build
build:
	go run ../../cmd/workers-assets-gen
	tinygo build -o ./build/app.wasm -target wasm -no-debug ./...

.PHONY: deploy
deploy:
	npx wrangler deploy
