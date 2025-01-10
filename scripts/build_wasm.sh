#!/bin/bash

# Ensure dist directory exists
mkdir -p dist

# Copy wasm_exec.js from Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dist/

# Create version.go with temporary build hash
mkdir -p frontend/internal
cat > frontend/internal/version.go << EOL
package internal

const BuildHash = "building"
EOL

# Build the frontend
GOOS=js GOARCH=wasm go build -o dist/main.wasm frontend/main.go 

# Generate hash of the wasm file
WASM_HASH=$(shasum -a 256 dist/main.wasm | cut -d' ' -f1 | head -c 8)

# Update version.go with the actual WASM hash
cat > frontend/internal/version.go << EOL
package internal

// BuildHash is automatically generated during build
const BuildHash = "${WASM_HASH}"
EOL

# Rebuild with the correct hash
GOOS=js GOARCH=wasm go build -o dist/main.wasm frontend/main.go 

# Copy index.html to dist if it exists
if [ -f frontend/index.html ]; then
    cp frontend/index.html dist/
fi

echo "Built WASM with hash: ${WASM_HASH}" 