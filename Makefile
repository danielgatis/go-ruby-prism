PRISM_SOURCE_DIR := prism
WASM_IMAGE := prism-wasm

.PHONY: all init_submodules wasm_build config generate format test clean

all: init_submodules wasm_build config generate format

init_submodules:
		@echo "Initializing submodules"
		git submodule update --init --recursive

# Builds prism.wasm inside a container so that neither Ruby nor the WASI SDK
# has to be installed on the host. prism.wasm is committed, so this only needs
# to run when the prism submodule moves.
# The image is a scratch layer holding only prism.wasm, so the module is
# copied out of a container rather than run out of one.
wasm_build:
		@echo "Building wasm"
		docker build -f Dockerfile.wasm -t $(WASM_IMAGE) .
		$(eval CID := $(shell docker create $(WASM_IMAGE)))
		docker cp $(CID):/prism.wasm wasm/prism.wasm
		docker rm -v $(CID) >/dev/null
		@echo "Wrote wasm/prism.wasm"

# Copies the node definitions out of the pinned submodule. The generator reads
# this file, so it has to stay in step with prism.wasm.
config:
		@echo "Syncing config.yml from $(PRISM_SOURCE_DIR)"
		cp $(PRISM_SOURCE_DIR)/config.yml config.yml

generate:
		@echo "Generating go files"
		go generate ./...
		go fmt ./...

format:
		@echo "Formatting go files"
		go fmt ./...

test:
		@echo "Running tests"
		go test -race ./...

clean:
		@echo "Cleaning up"
		rm -fr wasm/prism.wasm
		docker image rm -f $(WASM_IMAGE) 2>/dev/null || true
