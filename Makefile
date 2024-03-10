PROJECT_DIR=$(shell pwd)
PLUGIN_DIR=$(PROJECT_DIR)/pkg/plugins
DIST_DIR=$(PROJECT_DIR)/dist
DIST_PLUGIN_DIR=$(DIST_DIR)/plugins

DOCKER_DIR=$(PROJECT_DIR)/deployments
DOCKER_PROJECT=postmanq
DOCKER_COMPOSE=$(DOCKER_DIR)/docker-compose.yml
DOCKER_COMPOSE_INFRA=$(DOCKER_DIR)/docker-compose.infra.yml
DOCKER_COMPOSE_POSTMANQ=$(DOCKER_DIR)/docker-compose.postmanq.yml

PLUGINS=$(shell ls -d $(PLUGIN_DIR)/*)

GOCMD=go
GOBUILD=go build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

-include .env
export

deps:
	$(GOCMD) mod vendor

test:
	$(GOTEST) -race -parallel 8 ./...

lint:
	golangci-lint run --skip-files=module.go

clean:
	rm -rf $(DIST_DIR)

buf_update:
	buf mod update

buf_generate:
	buf generate
	cp -R $(DIST_DIR)/github.com/postmanq/postmanq/pkg $(PROJECT_DIR)
	rm -rf $(DIST_DIR)/github.com

build_plugins:
	$(foreach PLUGIN, $(PLUGINS), $(GOBUILD) -buildmode=plugin -o $(DIST_PLUGIN_DIR)/$(shell basename $(PLUGIN))-$(GOOS)-$(GOARCH) $(PLUGIN_DIR)/$(shell basename $(PLUGIN))/module.go;)

build_postmanq:
	$(GOBUILD) -o $(DIST_DIR)/postmanq-$(GOOS)-$(GOARCH) $(PROJECT_DIR)/cmd/postmanq/main.go

build: deps build_postmanq build_plugins

infra_up:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) up -d --force-recreate

infra_down:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_INFRA) down --remove-orphans

postmanq_up:
	docker-compose -p $(DOCKER_PROJECT) -f $(DOCKER_COMPOSE_POSTMANQ) up -d --build --force-recreate
