# Definition of binary names and build directory
AGENT_BINARY := agent
ORCHESTRATOR_BINARY := orchestrator
BUILD_DIR := build

# Definition of Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Definition of source files
AGENT_MAIN := cmd/agent/main.go
ORCHESTRATOR_MAIN := cmd/orchestrator/main.go

# Load variables from .env file if it exists
ifneq (,$(wildcard .env))
    include .env
endif

# Define environment variables with default values
export COMPUTING_POWER ?= $(or $(COMPUTING_POWER),4)
export TIME_ADDITION_MS ?= $(or $(TIME_ADDITION_MS),1000)
export TIME_SUBTRACTION_MS ?= $(or $(TIME_SUBTRACTION_MS),1000)
export TIME_MULTIPLICATIONS_MS ?= $(or $(TIME_MULTIPLICATIONS_MS),2000)
export TIME_DIVISIONS_MS ?= $(or $(TIME_DIVISIONS_MS),2000)
export ORCHESTRATOR_URL ?= http://localhost:8080
export PORT ?= 8080

# Create the build directory if it does not exist
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Definition of available commands
.PHONY: all build clean test deps run run-dev run-prod check-deps test-short test-coverage lint stop run-agent run-orchestrator run-race help

# Display the list of available commands
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

# Full build cycle: clean, install dependencies, build
all: clean deps build

# Download and verify dependencies
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Check dependencies and fix the module
check-deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

# Build binaries
build: deps $(BUILD_DIR)
	# Linux builds
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-linux-amd64 $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-linux-amd64 $(ORCHESTRATOR_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-linux-arm64 $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-linux-arm64 $(ORCHESTRATOR_MAIN)
	# Windows builds
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-windows-amd64.exe $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-windows-amd64.exe $(ORCHESTRATOR_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-windows-arm64.exe $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-windows-arm64.exe $(ORCHESTRATOR_MAIN)
	# macOS builds
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-macos-amd64 $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-macos-amd64 $(ORCHESTRATOR_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY)-macos-arm64 $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY)-macos-arm64 $(ORCHESTRATOR_MAIN)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests with race detection and code coverage
test:
	$(GOTEST) -v -race -cover ./...

test-short:
	$(GOTEST) -v -short ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run static code analyzer
lint:
	golangci-lint run ./...

# Macro for printing environment variables
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

# Run services
run: build
	@echo "Starting services..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

# Run in development mode with reduced delays
run-dev: export COMPUTING_POWER=2
run-dev: export TIME_ADDITION_MS=100
run-dev: export TIME_SUBTRACTION_MS=100
run-dev: export TIME_MULTIPLICATIONS_MS=200
run-dev: export TIME_DIVISIONS_MS=200
run-dev: build
	@echo "Starting services in development mode..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

# Run in production mode with increased power
run-prod: export COMPUTING_POWER=8
run-prod: build
	@echo "Starting services in production mode..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

# Stop services
stop:
	@if [ -f $(BUILD_DIR)/orchestrator.pid ]; then \
		kill $$(cat $(BUILD_DIR)/orchestrator.pid) || true; \
		rm $(BUILD_DIR)/orchestrator.pid; \
	fi
	@if [ -f $(BUILD_DIR)/agent.pid ]; then \
		kill $$(cat $(BUILD_DIR)/agent.pid) || true; \
		rm $(BUILD_DIR)/agent.pid; \
	fi

# Run only the agent
run-agent: build
	@echo "Starting agent..."
	$(print_env)
	@$(BUILD_DIR)/$(AGENT_BINARY)

# Run only the orchestrator
run-orchestrator: build
	@echo "Starting orchestrator..."
	@echo "Environment variables:"
	@echo "  PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)

# Build and run with race detection
run-race: clean
	$(GOBUILD) -race -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -race -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)
	@echo "Starting services with race detection..."
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid
