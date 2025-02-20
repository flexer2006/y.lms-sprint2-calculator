# Определение имен бинарников и каталога сборки
AGENT_BINARY := agent
ORCHESTRATOR_BINARY := orchestrator
BUILD_DIR := build

# Определение команд Go
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod

# Определение исходных файлов
AGENT_MAIN := cmd/agent/main.go
ORCHESTRATOR_MAIN := cmd/orchestrator/main.go

# Загрузка переменных из .env файла, если он существует
ifneq (,$(wildcard .env))
    include .env
endif

# Определение переменных окружения с значениями по умолчанию
export COMPUTING_POWER ?= $(or $(COMPUTING_POWER),4)
export TIME_ADDITION_MS ?= $(or $(TIME_ADDITION_MS),1000)
export TIME_SUBTRACTION_MS ?= $(or $(TIME_SUBTRACTION_MS),1000)
export TIME_MULTIPLICATIONS_MS ?= $(or $(TIME_MULTIPLICATIONS_MS),2000)
export TIME_DIVISIONS_MS ?= $(or $(TIME_DIVISIONS_MS),2000)
export ORCHESTRATOR_URL ?= http://localhost:8080
export PORT ?= 8080

# Создание каталога сборки, если он не существует
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Определение доступных команд
.PHONY: all build clean test deps run run-dev run-prod check-deps test-short test-coverage lint stop run-agent run-orchestrator run-race help

# Вывод списка доступных команд
help:
	@echo "Available targets:"
	@echo "  all              - clean, deps, build"
	@echo "  deps             - download and verify dependencies"
	@echo "  check-deps       - tidy and verify modules"
	@echo "  build            - build binaries (depends on $(BUILD_DIR))"
	@echo "  clean            - clean build artifacts"
	@echo "  test             - run tests with race and coverage"
	@echo "  test-short       - run short tests"
	@echo "  test-coverage    - generate coverage report"
	@echo "  lint             - run golangci-lint"
	@echo "  run              - run orchestrator and agent"
	@echo "  run-dev          - run in development mode"
	@echo "  run-prod         - run in production mode"
	@echo "  stop             - stop running services"
	@echo "  run-agent        - run agent only"
	@echo "  run-orchestrator - run orchestrator only"
	@echo "  run-race         - run services with race detection"

# Полный цикл сборки: очистка, установка зависимостей, сборка
all: clean deps build

# Загрузка и проверка зависимостей
deps:
	$(GOMOD) download
	$(GOMOD) verify

# Проверка зависимостей и исправление модуля
check-deps:
	$(GOCMD) mod tidy
	$(GOCMD) mod verify

# Сборка бинарников
build: deps $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)

# Очистка сборочных артефактов
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Запуск тестов с детекцией гонок и покрытием кода
test:
	$(GOTEST) -v -race -cover ./...

test-short:
	$(GOTEST) -v -short ./...

test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Запуск статического анализатора кода
lint:
	golangci-lint run ./...

# Макрос для вывода переменных окружения
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

# Запуск сервисов
run: build
	@echo "Starting services..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

# Запуск в режиме разработки с уменьшенными задержками
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

# Запуск в продакшн-режиме с увеличенной мощностью
run-prod: export COMPUTING_POWER=8
run-prod: build
	@echo "Starting services in production mode..."
	$(print_env)
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid

# Остановка сервисов
stop:
	@if [ -f $(BUILD_DIR)/orchestrator.pid ]; then \
		kill $$(cat $(BUILD_DIR)/orchestrator.pid) || true; \
		rm $(BUILD_DIR)/orchestrator.pid; \
	fi
	@if [ -f $(BUILD_DIR)/agent.pid ]; then \
		kill $$(cat $(BUILD_DIR)/agent.pid) || true; \
		rm $(BUILD_DIR)/agent.pid; \
	fi

# Запуск только агента
run-agent: build
	@echo "Starting agent..."
	$(print_env)
	@$(BUILD_DIR)/$(AGENT_BINARY)

# Запуск только оркестратора
run-orchestrator: build
	@echo "Starting orchestrator..."
	@echo "Environment variables:"
	@echo "  PORT: $(PORT)"
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY)

# Сборка и запуск с детекцией гонок
run-race: clean
	$(GOBUILD) -race -o $(BUILD_DIR)/$(AGENT_BINARY) $(AGENT_MAIN)
	$(GOBUILD) -race -o $(BUILD_DIR)/$(ORCHESTRATOR_BINARY) $(ORCHESTRATOR_MAIN)
	@echo "Starting services with race detection..."
	@$(BUILD_DIR)/$(ORCHESTRATOR_BINARY) & echo $$! > $(BUILD_DIR)/orchestrator.pid
	@$(BUILD_DIR)/$(AGENT_BINARY) & echo $$! > $(BUILD_DIR)/agent.pid
