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

		cd prism && bundle install
		cd prism && bundle exec rake compile

		wget https://github.com/WebAssembly/wasi-sdk/releases/download/wasi-sdk-25/wasi-sdk-25.0-arm64-macos.tar.gz
		tar xvf wasi-sdk-25.0-arm64-macos.tar.gz

		cd $(PRISM_SOURCE_DIR) && PRISM_SERIALIZE_ONLY_SEMANTICS_FIELDS=1 bundle exec rake templates
		cd $(PRISM_SOURCE_DIR) && make java-wasm WASI_SDK_PATH=../wasi-sdk-25.0-arm64-macos PRISM_SERIALIZE_ONLY_SEMANTICS_FIELDS=1
		cp -f prism/java-wasm/src/test/resources/prism.wasm wasm

clean:
		@echo "Cleaning up"
		rm -fr prism/java-wasm/src/test/resources/prism.wasm
		rm -fr wasm/prism.wasm
		rm -fr wasi-sdk-25.0-arm64-macos.tar.gz
		rm -fr wasi-sdk-25.0-arm64-macos
		cd prism && bundle exec rake clean
		cd $(PRISM_SOURCE_DIR) && make clean
