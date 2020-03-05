PROJECT_DIR=$(shell pwd)
MODULE_DIR=$(PROJECT_DIR)/module
OUT_DIR=$(PROJECT_DIR)/out
OUT_MODULE_DIR=$(OUT_DIR)/module

MODULES=$(shell ls -d $(MODULE_DIR)/*/)

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

test:
	$(GOTEST) -v ./...

clean:
	rm -rf $(OUT_DIR)

build-postmanq:
	$(GOBUILD) -o $(OUT_DIR)/postmanq $(PROJECT_DIR)/cmd/postmanq/postmanq.go

build-modules:
	mkdir -p $(OUT_MODULE_DIR)
	$(foreach MODULE, $(MODULES), $(GOBUILD) -buildmode=plugin -o $(OUT_MODULE_DIR)/$(shell basename $(MODULE))/module.so $(MODULE_DIR)/$(shell basename $(MODULE))/module.go;)

build: clean build-modules build-postmanq

run: test build
	$(OUT_DIR)/postmanq -c $(PROJECT_DIR)/config.yml