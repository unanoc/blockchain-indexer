#! /usr/bin/make -f

# ---------------------------------------------------------
# Project variables.
# ---------------------------------------------------------
PROJECT_NAME := $(shell basename "$(PWD)")
PACKAGE := github.com/unanoc/$(PROJECT_NAME)
VERSION := $(shell git describe --tags 2>/dev/null || git describe --all)
BUILD := $(shell git rev-parse --short HEAD)
DATETIME := $(shell date +"%Y.%m.%d-%H:%M:%S")

# ---------------------------------------------------------
# Use linker flags to provide version/build settings.
# ---------------------------------------------------------
LDFLAGS=-ldflags "-X=$(PACKAGE)/build.Version=$(VERSION) -X=$(PACKAGE)/build.Build=$(BUILD) -X=$(PACKAGE)/build.Date=$(DATETIME)"

# ---------------------------------------------------------
# Go related variables.
# ---------------------------------------------------------
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin

# ---------------------------------------------------------
# Go files.
# ---------------------------------------------------------
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# ---------------------------------------------------------
# Docker Variables
# ---------------------------------------------------------
DOCKER_COMPOSE_FILE_LOCAL_ENV ?= ./docker-compose.yaml

# ---------------------------------------------------------
# API
# ---------------------------------------------------------
api: go-build-api start-api

start-api:
	@echo "  >  Starting api"
	PROMETHEUS_SUBSYSTEM=api $(GOBIN)/api

go-build-api:
	@echo "  >  Building api binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/api ./cmd/api

# ---------------------------------------------------------
# Blockproducer
# ---------------------------------------------------------
blockproducer: go-build-blockproducer start-blockproducer

start-blockproducer:
	@echo "  >  Starting blockproducer"
	PROMETHEUS_SUBSYSTEM=blockproducer $(GOBIN)/blockproducer

go-build-blockproducer:
	@echo "  >  Building blockproducer binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/blockproducer ./cmd/blockproducer

# ---------------------------------------------------------
# Blockconsumer
# ---------------------------------------------------------
blockconsumer: go-build-blockconsumer start-blockconsumer

start-blockconsumer:
	@echo "  >  Starting blockconsumer"
	PROMETHEUS_SUBSYSTEM=blockconsumer $(GOBIN)/blockconsumer

go-build-blockconsumer:
	@echo "  >  Building blockconsumer binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/blockconsumer ./cmd/blockconsumer

# ---------------------------------------------------------
# Nodes
# ---------------------------------------------------------
nodes: go-build-nodes start-nodes

start-nodes:
	@echo "  >  Starting nodes"
	PROMETHEUS_SUBSYSTEM=nodes $(GOBIN)/nodes

go-build-nodes:
	@echo "  >  Building nodes binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/nodes ./cmd/nodes

# ---------------------------------------------------------
# Transactionconsumer
# ---------------------------------------------------------
transactionconsumer: go-build-transactionconsumer start-transactionconsumer

start-transactionconsumer:
	@echo "  >  Starting transactionconsumer"
	PROMETHEUS_SUBSYSTEM=transactionconsumer $(GOBIN)/transactionconsumer

go-build-transactionconsumer:
	@echo "  >  Building transactionconsumer binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/transactionconsumer ./cmd/transactionconsumer

# ---------------------------------------------------------
# Code Checking
# ---------------------------------------------------------
check: fmt lint unit-test swag

unit-test:
	@echo "  >  Running unit tests"
	go clean -testcache
	GOBIN=$(GOBIN) go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...

fmt:
	@echo "  >  Format all go files"
	GOBIN=$(GOBIN) gofmt -w ${GOFMT_FILES}

lint-install:
ifeq ("$(wildcard bin/golangci-lint)","")
	@echo "  >  Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

lint: lint-install
	@echo "  >  Running golint"
	bin/golangci-lint run --timeout=3m

swag:
	@echo "  >  Generating swagger files"
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g cmd/api/main.go

# ---------------------------------------------------------
# Local Environment
# ---------------------------------------------------------
up:
	@echo "  >  Run local environment"
	docker-compose -f ${DOCKER_COMPOSE_FILE_LOCAL_ENV} up -d

# ---------------------------------------------------------
# Code generation
# ---------------------------------------------------------
generate-coins:
	@echo "  >  Generating coin file"
	GOBIN=$(GOBIN) go run -tags=coins pkg/primitives/coin/gen.go
