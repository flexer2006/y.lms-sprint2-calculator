# Binary names
AGENT_BINARY=agent
ORCHESTRATOR_BINARY=orchestrator
BUILD_DIR=build

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Source files
AGENT_MAIN=cmd/agent/main.go
ORCHESTRATOR_MAIN=cmd/orchestrator/main.go

# Default environment variables
export COMPUTING_POWER ?= 4
export TIME_ADDITION_MS ?= 1000
export TIME_SUBTRACTION_MS ?= 1000
export TIME_MULTIPLICATIONS_MS ?= 2000
export TIME_DIVISIONS_MS ?= 2000
export ORCHESTRATOR_URL ?= http://localhost:8080
export PORT ?= 8080

# Create build directory if it doesn't exist
$(shell mkdir -p $(BUILD_DIR))

.PHONY: all build clean test deps run run-dev run-prod check-deps test-short test-coverage lint

all: clean deps build

deps:
	$(GOMOD) download
	$(GOMOD) verify

check-deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

build: deps
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)

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

run: build
	@echo "Starting services..."
	@echo "Environment variables:"
	@echo "COMPUTING_POWER: $(COMPUTING_POWER)"
	@echo "TIME_ADDITION_MS: $(TIME_ADDITION_MS)"
	@echo "TIME_SUBTRACTION_MS: $(TIME_SUBTRACTION_MS)"
	@echo "TIME_MULTIPLICATIONS_MS: $(TIME_MULTIPLICATIONS_MS)"
	@echo "TIME_DIVISIONS_MS: $(TIME_DIVISIONS_MS)"
	@echo "ORCHESTRATOR_URL: $(ORCHESTRATOR_URL)"
	@echo "PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

run-dev: export COMPUTING_POWER=2
run-dev: export TIME_ADDITION_MS=100
run-dev: export TIME_SUBTRACTION_MS=100
run-dev: export TIME_MULTIPLICATIONS_MS=200
run-dev: export TIME_DIVISIONS_MS=200
run-dev: build
	@echo "Starting services in development mode..."
	@echo "Environment variables:"
	@echo "COMPUTING_POWER: $(COMPUTING_POWER)"
	@echo "TIME_ADDITION_MS: $(TIME_ADDITION_MS)"
	@echo "TIME_SUBTRACTION_MS: $(TIME_SUBTRACTION_MS)"
	@echo "TIME_MULTIPLICATIONS_MS: $(TIME_MULTIPLICATIONS_MS)"
	@echo "TIME_DIVISIONS_MS: $(TIME_DIVISIONS_MS)"
	@echo "ORCHESTRATOR_URL: $(ORCHESTRATOR_URL)"
	@echo "PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

run-prod: export COMPUTING_POWER=8
run-prod: export TIME_ADDITION_MS=1000
run-prod: export TIME_SUBTRACTION_MS=1000
run-prod: export TIME_MULTIPLICATIONS_MS=2000
run-prod: export TIME_DIVISIONS_MS=2000
run-prod: build
	@echo "Starting services in production mode..."
	@echo "Environment variables:"
	@echo "COMPUTING_POWER: $(COMPUTING_POWER)"
	@echo "TIME_ADDITION_MS: $(TIME_ADDITION_MS)"
	@echo "TIME_SUBTRACTION_MS: $(TIME_SUBTRACTION_MS)"
	@echo "TIME_MULTIPLICATIONS_MS: $(TIME_MULTIPLICATIONS_MS)"
	@echo "TIME_DIVISIONS_MS: $(TIME_DIVISIONS_MS)"
	@echo "ORCHESTRATOR_URL: $(ORCHESTRATOR_URL)"
	@echo "PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

stop:
	@if [ -f $(BUILD_DIR)/orchestrator.pid ]; then \
		kill `cat $(BUILD_DIR)/orchestrator.pid` || true; \
		rm $(BUILD_DIR)/orchestrator.pid; \
	fi
	@if [ -f $(BUILD_DIR)/agent.pid ]; then \
		kill `cat $(BUILD_DIR)/agent.pid` || true; \
		rm $(BUILD_DIR)/agent.pid; \
	fi

run-agent: build
	@echo "Starting agent..."
	@echo "Environment variables:"
	@echo "COMPUTING_POWER: $(COMPUTING_POWER)"
	@echo "TIME_ADDITION_MS: $(TIME_ADDITION_MS)"
	@echo "TIME_SUBTRACTION_MS: $(TIME_SUBTRACTION_MS)"
	@echo "TIME_MULTIPLICATIONS_MS: $(TIME_MULTIPLICATIONS_MS)"
	@echo "TIME_DIVISIONS_MS: $(TIME_DIVISIONS_MS)"
	@echo "ORCHESTRATOR_URL: $(ORCHESTRATOR_URL)"
	$(BUILD_DIR)/$(AGENT_BINARY)

run-orchestrator: build
	@echo "Starting orchestrator..."
	@echo "Environment variables:"
	@echo "PORT: $(PORT)"
	$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)

run-race: clean
	$(GOBUILD) -race -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -race -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)
	@echo "Starting services with race detection..."
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid