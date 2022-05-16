GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)

GOCMD=go
GOBUILD=$(GOCMD) build
BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist

API_PROTO_FILES=$(shell find api -name *.proto)
ERROR_PROTO_FILES=$(shell find api/v1/errors -name *.proto)
CONFIG_PROTO_FILES=$(shell find app/pkg/config -name *.proto)

KUBEOPS_SRC=$(BASEPATH)/cmd
KUBEOPS_SERVER_NAME=apiserver
KUBEOPS_INVENTORY_NAME=inventory
KUBEOPS_CLIENT_NAME=opsctl

BIN_DIR=usr/local/bin
CONFIG_DIR=etc/kubeops
BASE_DIR=var/kubeops

.PHONY: init
# init env
init:
	go get -u github.com/go-kratos/kratos/cmd/kratos/v2
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2
	go get -u github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2
	go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
	go get -u github.com/google/wire/cmd/wire

.PHONY: build_linux
build_linux:
	GOOS=linux  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_SERVER_NAME) $(KUBEOPS_SRC)/server
	GOOS=linux  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_INVENTORY_NAME) $(KUBEOPS_SRC)/inventory
	GOOS=linux  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_CLIENT_NAME) $(KUBEOPS_SRC)/opsctl
	mkdir -p $(BUILDDIR)/$(CONFIG_DIR) && cp -r  $(BASEPATH)/conf/* $(BUILDDIR)/$(CONFIG_DIR)
	mkdir -p $(BUILDDIR)/$(BASE_DIR)/plugins/callback && cp  $(BASEPATH)/plugin/* $(BUILDDIR)/$(BASE_DIR)/plugins/callback

.PHONY: build_darwin
build_darwin:
	GOOS=darwin  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_SERVER_NAME) $(KUBEOPS_SRC)/server
	GOOS=darwin  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_INVENTORY_NAME) $(KUBEOPS_SRC)/inventory
	GOOS=darwin  GOARCH=$(GOARCH) $(GOBUILD) -o $(BUILDDIR)/$(BIN_DIR)/$(KUBEOPS_CLIENT_NAME) $(KUBEOPS_SRC)/opsctl
	mkdir -p $(BUILDDIR)/$(CONFIG_DIR) && cp -r  $(BASEPATH)/conf/* $(BUILDDIR)/$(CONFIG_DIR)
	mkdir -p $(BUILDDIR)/$(BASE_DIR)/plugins/callback && cp  $(BASEPATH)/plugin/* $(BUILDDIR)/$(BASE_DIR)/plugins/callback

.PHONY: clean
clean:
	rm -rf $(BUILDDIR)

.PHONY: docker
docker:
	@echo "build docker images"
	docker build -t pipperman/kubeops:v0.0.1 --build-arg GOARCH=$(GOARCH) .

.PHONY: grpc
# generate grpc code
grpc:
	protoc --proto_path=. \
		--proto_path=./third_party \
		--go_out=paths=source_relative:. \
		--go-grpc_out=paths=source_relative:. \
		$(API_PROTO_FILES)

.PHONY: http
# generate http code
http:
	protoc --proto_path=. \
		--proto_path=./third_party \
		--go_out=paths=source_relative:. \
		--go-http_out=paths=source_relative:. \
		$(API_PROTO_FILES)

.PHONY: errors
# generate errors code
errors:
	protoc --proto_path=. \
        --proto_path=./third_party \
        --go_out=paths=source_relative:. \
        --go-errors_out=paths=source_relative:. \
        $(ERROR_PROTO_FILES)

.PHONY: struct
# generate code struct
struct:
	protoc --proto_path=. \
       --proto_path=./third_party \
       --go_out=paths=source_relative:. \
       $(CONFIG_PROTO_FILES)

.PHONY: swagger
# generate swagger
swagger:
	protoc --proto_path=. \
	    --proto_path=./third_party \
	    --openapiv2_out . \
	    --openapiv2_opt logtostderr=true \
        $(API_PROTO_FILES)

.PHONY: wire
# generate wire
wire:
	cd cmd/server && wire

.PHONY: generate
# generate client code
generate:
	go generate ./...

.PHONY: all
# generate all
all:
	make grpc;
	make http;
	make struct;
	make errors;
	make swagger;
	make wire;