#include .env

PROJECT_DIR=$(shell pwd)
PLUGIN_DIR=$(PROJECT_DIR)/pkg/plugins
OUT_DIR=$(PROJECT_DIR)/out
OUT_PLUGIN_DIR=$(OUT_DIR)/plugins
MOCK_DIR=$(PROJECT_DIR)/mock

DOCKER_DIR=$(PROJECT_DIR)/deployments
DOCKER_PROJECT=postmanq
DOCKER_COMPOSE=$(DOCKER_DIR)/docker-compose.yml
DOCKER_COMPOSE_INFRA=$(DOCKER_DIR)/docker-compose.infra.yml

PLUGINS=$(shell ls -d $(PLUGIN_DIR)/* | cut -d '/' -f 3)

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BUILD_PROTO_DIR=$(PROJECT_DIR)/api/proto/postmanq
BUILD_GEN_DIR=$(PROJECT_DIR)/pkg/gen

-include .env
export

deps:
	$(GOCMD) mod vendor

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run --skip-files=module.go

clean:
	rm -rf $(OUT_DIR)
	rm -rf $(MOCK_DIR)

buf_update:
	buf mod update

buf_generate:
	rm -rf $(BUILD_GEN_DIR)
	mkdir -p $(BUILD_GEN_DIR)
	buf generate --path $(BUILD_PROTO_DIR)

build_postmanq:
	$(GOBUILD) -o $(OUT_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq.go

build_plugins:
	rm -rf $(OUT_PLUGIN_DIR)
	mkdir -p $(OUT_PLUGIN_DIR)
	$(foreach PLUGIN, $(PLUGINS), $(GOBUILD) -buildmode=plugin -o $(OUT_PLUGIN_DIR)/$(PLUGIN).so $(PLUGIN_DIR)/$(PLUGIN)/plugin/main.go;)

build: clean create-mocks test lint build-modules build-postmanq

run: build
	$(OUT_DIR)/postmanq -c $(PROJECT_DIR)/config.yml

infra_up:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) up -d --force-recreate

infra_down:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) down --remove-orphans