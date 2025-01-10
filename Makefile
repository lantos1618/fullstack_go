.PHONY: build_frontend build_backend run clean kill-server

build_frontend:
	GOOS=js GOARCH=wasm go build -o dist/main.wasm frontend/main.go
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js dist/
	cp frontend/index.html dist/

build_backend:
	go build -o tmp/main .

run_backend: build_backend
	./tmp/main

dev: kill-server
	air

clean:
	rm -rf dist/* tmp/*

kill-server:
	@lsof -ti :8080 | xargs kill -9 2>/dev/null || true 