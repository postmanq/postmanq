PROJECT_DIR=$(shell pwd)
PLUGIN_DIR=$(PROJECT_DIR)/pkg/plugins
DIST_DIR=$(PROJECT_DIR)/dist
DIST_PLUGIN_DIR=$(DIST_DIR)/plugins
MOCK_DIR=$(PROJECT_DIR)/mock

DOCKER_DIR=$(PROJECT_DIR)/deployments
DOCKER_PROJECT=postmanq
DOCKER_COMPOSE=$(DOCKER_DIR)/docker-compose.yml
DOCKER_COMPOSE_INFRA=$(DOCKER_DIR)/docker-compose.infra.yml

PLUGINS=$(shell ls -d $(PLUGIN_DIR)/*)

GOCMD=go
GOBUILD=$(GOCMD) build -ldflags="-extldflags=-Wl,-ld_classic"
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

BUILD_PROTO_DIR=$(PROJECT_DIR)/api/proto/postmanq
BUILD_GEN_DIR=$(PROJECT_DIR)/pkg/common/gen

-include .env
export

deps:
	$(GOCMD) mod vendor

test:
	$(GOTEST) -v ./...

lint:
	golangci-lint run --skip-files=module.go

clean:
	rm -rf $(DIST_DIR)
	rm -rf $(MOCK_DIR)

buf_update:
	buf mod update

buf_generate:
	rm -rf $(BUILD_GEN_DIR)
	mkdir -p $(BUILD_GEN_DIR)
	buf generate --path $(BUILD_PROTO_DIR)

build_plugins:
	rm -rf $(DIST_PLUGIN_DIR)
	mkdir -p $(DIST_PLUGIN_DIR)
	$(foreach PLUGIN, $(PLUGINS), $(GOBUILD) -buildmode=plugin -o $(DIST_PLUGIN_DIR)/$(shell basename $(PLUGIN)).so $(PLUGIN_DIR)/$(shell basename $(PLUGIN))/module.go;)

build_postmanq:
	$(GOBUILD) -o $(DIST_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq/main.go

run: build_plugins build_postmanq
	$(DIST_DIR)/postmanq -c $(PROJECT_DIR)/config.yml

infra_up:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) up -d --force-recreate

infra_down:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) down --remove-orphans