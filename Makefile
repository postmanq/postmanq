include .env

PROJECT_DIR=$(shell pwd)
MODULE_DIR=$(PROJECT_DIR)/module
OUT_DIR=$(PROJECT_DIR)/out
OUT_MODULE_DIR=$(OUT_DIR)/module
MOCK_DIR=$(PROJECT_DIR)/mock
MOCK_MODULE_DIR=$(PROJECT_DIR)/mock/module
MOCK_COMPONENT_INTERFACES=InitComponent ReceiveComponent SendComponent ProcessComponent

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

lint:
	golangci-lint run --skip-files=module.go

clean:
	rm -rf $(OUT_DIR)
	rm -rf $(MOCK_DIR)

build-postmanq:
	$(GOBUILD) -o $(OUT_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq.go

build-modules:
	mkdir -p $(OUT_MODULE_DIR)
	$(foreach MODULE, $(MODULES), $(GOBUILD) -buildmode=plugin -o $(OUT_MODULE_DIR)/$(shell basename $(MODULE))/module.so $(MODULE_DIR)/$(shell basename $(MODULE))/module.go;)

build: clean create-mocks test lint build-modules build-postmanq

run: build
	$(OUT_DIR)/postmanq -c $(PROJECT_DIR)/config.yml

create-mocks:
	$(foreach MOCK_COMPONENT_INTERFACE, $(MOCK_COMPONENT_INTERFACES), $(MOCKERY) -dir=$(MODULE_DIR) -output=$(MOCK_MODULE_DIR) -outpkg=module -case=snake -name $(MOCK_COMPONENT_INTERFACE);)
	$(foreach MOCK_PACKAGE, $(MOCK_PACKAGES), $(MOCKERY) -dir=$(MODULE_DIR)/$(MOCK_PACKAGE) -output=$(MOCK_MODULE_DIR)/$(MOCK_PACKAGE) -outpkg=$(shell basename $(MOCK_PACKAGE)) -all -case=snake;)