#include .env

PROJECT_DIR=$(shell pwd)
PLUGIN_DIR=$(PROJECT_DIR)/pkg/plugins
OUT_DIR=$(PROJECT_DIR)/out
OUT_PLUGIN_DIR=$(OUT_DIR)/plugins
MOCK_DIR=$(PROJECT_DIR)/mock

PLUGINS=$(shell ls -d $(PLUGIN_DIR)/* | cut -d '/' -f 3)

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

deps:
	$(GOCMD) mod vendor

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run --skip-files=module.go

clean:
	rm -rf $(OUT_DIR)
	rm -rf $(MOCK_DIR)

build-postmanq:
	$(GOBUILD) -o $(OUT_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq.go

build-plugins:
	rm -rf $(OUT_PLUGIN_DIR)
	mkdir -p $(OUT_PLUGIN_DIR)
	$(foreach PLUGIN, $(PLUGINS), $(GOBUILD) -buildmode=plugin -o $(OUT_PLUGIN_DIR)/$(PLUGIN).so $(PLUGIN_DIR)/$(PLUGIN)/plugin/main.go;)

build: clean create-mocks test lint build-modules build-postmanq

run: build
	$(OUT_DIR)/postmanq -c $(PROJECT_DIR)/config.yml