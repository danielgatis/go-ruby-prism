PRISM_SOURCE_DIR := prism
ROOT_SOURCE_DIR := .

.PHONY: all

all: copy_config copy_template generate

copy_config:
		@echo "Copying config.yml"
		@cp $(PRISM_SOURCE_DIR)/config.yml $(ROOT_SOURCE_DIR)/config.yml

copy_template:
		@echo "Copying template.rb"
		@cp $(PRISM_SOURCE_DIR)/templates/template.rb $(ROOT_SOURCE_DIR)/templates/template.rb

generate:
		@echo "Generating go files"
		@go generate ./...
		@go fmt ./...
