.PHONY: build watch clean


build-wasm:
	@echo "Building WASM..."
	@./scripts/build_wasm.sh

watch:
	@echo "Watching for changes..."
	@air

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist
	@rm -f frontend/internal/version.go

dev: clean build watch 