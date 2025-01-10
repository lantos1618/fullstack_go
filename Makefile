.PHONY: dev build clean

dev:
	@if ! command -v air > /dev/null; then \
		go install github.com/air-verse/air@latest; \
	fi
	$(shell go env GOPATH)/bin/air

build:
	@mkdir -p frontend
	@cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" frontend/
	@GOOS=js GOARCH=wasm go build -o frontend/main.wasm frontend/main.go
	@go build -o tmp/main .

clean:
	@rm -rf tmp
	@rm -f frontend/main.wasm
	@rm -f frontend/wasm_exec.js 