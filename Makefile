# Binary names and build directory
AGENT_BINARY := agent
ORCHESTRATOR_BINARY := orchestrator
BUILD_DIR := build

# Detect OS and architecture
ifeq ($(OS),Windows_NT)
    BINARY_SUFFIX := .exe
    DETECTED_OS := windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        DETECTED_OS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        DETECTED_OS := macos
    endif
    BINARY_SUFFIX :=
endif

ifeq ($(shell uname -m),x86_64)
    DETECTED_ARCH := amd64
else ifeq ($(shell uname -m),aarch64)
    DETECTED_ARCH := arm64
else
    DETECTED_ARCH := amd64
endif

# Go commands
GOCMD    := go
GOBUILD  := $(GOCMD) build
GOCLEAN  := $(GOCMD) clean
GOTEST   := $(GOCMD) test
GOGET    := $(GOCMD) get
GOMOD    := $(GOCMD) mod

# Source files
AGENT_MAIN         := cmd/agent/main.go
ORCHESTRATOR_MAIN  := cmd/orchestrator/main.go

# Load variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
endif

# Define environment variables with default values (if not set in .env)
export COMPUTING_POWER          ?= $(or $(COMPUTING_POWER),4)
export TIME_ADDITION_MS         ?= $(or $(TIME_ADDITION_MS),1000)
export TIME_SUBTRACTION_MS      ?= $(or $(TIME_SUBTRACTION_MS),1000)
export TIME_MULTIPLICATIONS_MS  ?= $(or $(TIME_MULTIPLICATIONS_MS),2000)
export TIME_DIVISIONS_MS        ?= $(or $(TIME_DIVISIONS_MS),2000)
export ORCHESTRATOR_URL         ?= $(or $(ORCHESTRATOR_URL),http://localhost:8080)
export PORT                     ?= $(or $(PORT),8080)

# New variables for CGO configuration:
# To compile with gcc (CGO enabled), call make with ENABLE_CGO=1 (by default 0)
ENABLE_CGO ?= 0
# CC can be overridden if needed; by default use gcc
CC ?= gcc

# Create the build directory if it doesn't exist
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Available commands
.PHONY: all build clean test deps run run-dev run-prod check-deps test-short test-coverage lint stop run-agent run-orchestrator run-race help docker-build docker-run docker-stop docker-clean docker-dev docker-prod

help:
	@echo "Available targets:"
	@echo "  all              - clean, deps, build"
	@echo "  deps             - download and verify dependencies"
	@echo "  check-deps       - tidy and verify modules"
	@echo "  build            - build binaries"
	@echo "  clean            - clean build artifacts"
	@echo "  test             - run tests with race and coverage"
	@echo "  test-short       - run short tests"
	@echo "  test-coverage    - generate coverage report"
	@echo "  lint             - run golangci-lint"
	@echo "  run              - run with .env variables"
	@echo "  run-dev          - run in development mode"
	@echo "  run-prod         - run in production mode"
	@echo "  stop             - stop running services"
	@echo "  run-agent        - run agent only"
	@echo "  run-orchestrator - run orchestrator only"
	@echo "  run-race         - run services with race detection"
	@echo "  docker-build     - build Docker images"
	@echo "  docker-run       - run Docker containers"
	@echo "  docker-stop      - stop Docker containers"
	@echo "  docker-clean     - clean Docker resources"
	@echo "  docker-dev       - run Docker containers in development mode"
	@echo "  docker-prod      - run Docker containers in production mode"

all: clean deps build

deps:
	$(GOMOD) download
	$(GOMOD) verify

check-deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

# Build binaries for current platform.
# Build commands now use the CGO configuration.
build: deps $(BUILD_DIR)
	@echo "Building with CGO_ENABLED=$(ENABLE_CGO) and CC=$(CC)"
	CGO_ENABLED=$(ENABLE_CGO) CC=$(CC) $(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-$(DETECTED_OS)-$(DETECTED_ARCH)$(BINARY_SUFFIX) $(AGENT_MAIN)
	CGO_ENABLED=$(ENABLE_CGO) CC=$(CC) $(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-$(DETECTED_OS)-$(DETECTED_ARCH)$(BINARY_SUFFIX) $(ORCHESTRATOR_MAIN)
	# Create symlinks or copies for generic names
	cd $(BUILD_DIR) && cp $(AGENT_BINARY)-$(DETECTED_OS)-$(DETECTED_ARCH)$(BINARY_SUFFIX) $(AGENT_BINARY)$(BINARY_SUFFIX)
	cd $(BUILD_DIR) && cp $(ORCHESTRATOR_BINARY)-$(DETECTED_OS)-$(DETECTED_ARCH)$(BINARY_SUFFIX) $(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test:
	$(GOTEST) -v -race -cover ./...

test-short:
	$(GOTEST) -v -short ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

define print_env
	@echo "Environment variables:"
	@echo "  COMPUTING_POWER: $(COMPUTING_POWER)"
	@echo "  TIME_ADDITION_MS: $(TIME_ADDITION_MS)"
	@echo "  TIME_SUBTRACTION_MS: $(TIME_SUBTRACTION_MS)"
	@echo "  TIME_MULTIPLICATIONS_MS: $(TIME_MULTIPLICATIONS_MS)"
	@echo "  TIME_DIVISIONS_MS: $(TIME_DIVISIONS_MS)"
	@echo "  ORCHESTRATOR_URL: $(ORCHESTRATOR_URL)"
	@echo "  PORT: $(PORT)"
endef

run: build
	@echo "Starting services..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/agent.pid

run-dev: export COMPUTING_POWER=2
run-dev: export TIME_ADDITION_MS=100
run-dev: export TIME_SUBTRACTION_MS=100
run-dev: export TIME_MULTIPLICATIONS_MS=200
run-dev: export TIME_DIVISIONS_MS=200
run-dev: build
	@echo "Starting services in development mode..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/agent.pid

run-prod: export COMPUTING_POWER=8
run-prod: build
	@echo "Starting services in production mode..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/agent.pid

stop:
	@if [ -f $(BUILD_DIR)/orchestrator.pid ]; then \
		kill $$(cat $(BUILD_DIR)/orchestrator.pid) || true; \
		rm $(BUILD_DIR)/orchestrator.pid; \
	fi
	@if [ -f $(BUILD_DIR)/agent.pid ]; then \
		kill $$(cat $(BUILD_DIR)/agent.pid) || true; \
		rm $(BUILD_DIR)/agent.pid; \
	fi

run-agent: build
	@echo "Starting agent..."
	$(print_env)
	@$(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX)

run-orchestrator: build
	@echo "Starting orchestrator..."
	@echo "Environment variables:"
	@echo "  PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX)

run-race: clean
	CGO_ENABLED=$(ENABLE_CGO) CC=$(CC) $(GOBUILD) -race -o $(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX) $(AGENT_MAIN)
	CGO_ENABLED=$(ENABLE_CGO) CC=$(CC) $(GOBUILD) -race -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX) $(ORCHESTRATOR_MAIN)
	@echo "Starting services with race detection..."
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY)$(BINARY_SUFFIX) & echo $$! > $(BUILD_DIR)/agent.pid

# Docker commands
docker-build:
	docker-compose build

docker-run:
	docker-compose up -d

docker-stop:
	docker-compose down

docker-clean:
	docker-compose down --rmi all --volumes --remove-orphans

docker-dev: export COMPUTING_POWER=2
docker-dev: export TIME_ADDITION_MS=100
docker-dev: export TIME_SUBTRACTION_MS=100
docker-dev: export TIME_MULTIPLICATIONS_MS=200
docker-dev: export TIME_DIVISIONS_MS=200
docker-dev:
	docker-compose up -d

docker-prod: export COMPUTING_POWER=8
docker-prod:
	docker-compose up -d