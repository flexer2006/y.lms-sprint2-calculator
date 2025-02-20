
AGENT_BINARY=agent
ORCHESTRATOR_BINARY=orchestrator


BUILD_DIR=build


GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod


AGENT_MAIN=cmd/agent/main.go
ORCHESTRATOR_MAIN=cmd/orchestrator/main.go


$(shell mkdir -p $(BUILD_DIR))

.PHONY: all build clean test deps run


all: clean deps build


deps:
	$(GOMOD) download
	$(GOMOD) verify


build: deps
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)


clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)


test:
	$(GOTEST) -v ./...


run: build
	@echo "Starting services..."
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
	$(BUILD_DIR)/$(AGENT_BINARY)


run-orchestrator: build
	@echo "Starting orchestrator..."
	$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)


run-race: clean
	$(GOBUILD) -race -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -race -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)
	@echo "Starting services with race detection..."
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid


help:
	@echo "Available commands:"
	@echo "  make all              - Clean, download dependencies and build all services"
	@echo "  make build           - Build both services"
	@echo "  make clean           - Clean build directory"
	@echo "  make deps            - Download and verify dependencies"
	@echo "  make test            - Run tests"
	@echo "  make run             - Build and run both services"
	@echo "  make stop            - Stop running services"
	@echo "  make run-agent       - Build and run agent service only"
	@echo "  make run-orchestrator - Build and run orchestrator service only"
	@echo "  make run-race        - Build and run services with race detection"
