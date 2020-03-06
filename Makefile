include .env

PROJECT_DIR=$(shell pwd)
MODULE_DIR=$(PROJECT_DIR)/module
OUT_DIR=$(PROJECT_DIR)/out
OUT_MODULE_DIR=$(OUT_DIR)/module
MOCK_DIR=$(PROJECT_DIR)/mock

MODULES=$(shell ls -d $(MODULE_DIR)/*/)

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
MOCKERY=$(GOBIN)/mockery

deps:
	$(GOCMD) mod vendor

test:
	$(GOTEST) -v ./...

clean:
	rm -rf $(OUT_DIR)
	rm -rf $(MOCK_DIR)

build-postmanq:
	$(GOBUILD) -o $(OUT_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq.go

build-modules:
	mkdir -p $(OUT_MODULE_DIR)
	$(foreach MODULE, $(MODULES), $(GOBUILD) -buildmode=plugin -o $(OUT_MODULE_DIR)/$(shell basename $(MODULE))/module.so $(MODULE_DIR)/$(shell basename $(MODULE))/module.go;)

build: clean create-mocks test build-modules build-postmanq

run: build
	$(OUT_DIR)/postmanq -c $(PROJECT_DIR)/config.yml

create-mocks:
	$(foreach MOCK_PACKAGE, $(MOCK_PACKAGES), $(MOCKERY) -dir=$(MODULE_DIR)/$(MOCK_PACKAGE) -output=$(MOCK_DIR)/$(MOCK_PACKAGE) -outpkg=$(shell basename $(MOCK_PACKAGE)) -all -case=snake;)