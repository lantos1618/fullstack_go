#!/bin/bash

# Copy wasm_exec.js from Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" frontend/

# Build the frontend
GOOS=js GOARCH=wasm go build -o frontend/main.wasm frontend/main.go 