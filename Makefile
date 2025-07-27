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
		go fmt ./...

format:
		@echo "Formatting go files"
		go fmt ./...

wasm_build:
		@echo "Building wasm"
		rm -fr prism/javascript/src/prism.wasm

		cd prism && bundle install
		cd prism && bundle exec rake compile

		if [ ! -d wasi-sdk-25.0-arm64-macos ]; then \
			if [ ! -f wasi-sdk-25.0-arm64-macos.tar.gz ]; then \
				wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-25/wasi-sdk-25.0-arm64-macos.tar.gz; \
			fi; \
			tar xvf wasi-sdk-25.0-arm64-macos.tar.gz; \
		fi

		cd $(PRISM_SOURCE_DIR) && bundle exec rake templates
		cd $(PRISM_SOURCE_DIR) && make wasm WASI_SDK_PATH=../wasi-sdk-25.0-arm64-macos
		cp -f prism/javascript/src/prism.wasm wasm

clean:
		@echo "Cleaning up"
		rm -fr priprism/javascript/src/prism.wasm
		rm -fr wasm/prism.wasm
		rm -fr wasi-sdk-25.0-arm64-macos.tar.gz
		rm -fr wasi-sdk-25.0-arm64-macos
		cd prism && bundle exec rake clean
		cd $(PRISM_SOURCE_DIR) && make clean
