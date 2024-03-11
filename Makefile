PRISM_SOURCE_DIR := prism
ROOT_SOURCE_DIR := .

.PHONY: all

all: init_submodules wasm_build generate format

init_submodules:
		@echo "Initializing submodules"
		git submodule update --init --recursive

generate:
		@echo "Generating go files"
		go generate ./...

format:
		@echo "Formatting go files"
		go fmt ./...

wasm_build:
		@echo "Building wasm"
		rm -fr prism/java-wasm/src/test/resources/prism.wasm
		cd $(PRISM_SOURCE_DIR) && make java-wasm
		cp -f prism/java-wasm/src/test/resources/prism.wasm wasm
