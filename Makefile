.PHONY: dev build clean

dev:
	@if ! command -v air > /dev/null; then \
		go install github.com/air-verse/air@latest; \
	fi
	$(shell go env GOPATH)/bin/air

build_frontend:
	@mkdir -p dist
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/
	@cp frontend/index.html dist/
	@GOOS=js GOARCH=wasm go build -o dist/main.wasm frontend/main.go

build_backend:
	@go build -o tmp/main .

clean:
	@rm -rf tmp dist 