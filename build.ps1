#!/usr/bin/env pwsh

param(
    [string]$Command = "run-dev"  # Default command argument
)

# Load environment variables from .env file
$envFile = ".env"
$DEFAULT_ENV = @{
    "ORCHESTRATOR_URL" = "http://localhost:8080"
    "PORT" = "8080"
}

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $DEFAULT_ENV[$matches[1].Trim()] = $matches[2].Trim()
        }
    }
}

foreach ($key in $DEFAULT_ENV.Keys) {
    if (-not (Test-Path env:$key)) {
        Set-Item -Path env:$key -Value $DEFAULT_ENV[$key]
    }
}

# Directories
$BUILD_DIR = "build"
if (-not (Test-Path $BUILD_DIR)) { New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null }

# Utility Functions
function Show-Help {
    Write-Host "Usage: .\build.ps1 <command>"
    Write-Host "Commands:"
    Write-Host "  build            - Build binaries"
    Write-Host "  clean            - Clean build artifacts"
    Write-Host "  test             - Run tests with coverage"
    Write-Host "  run-dev          - Run in development mode"
    Write-Host "  run-prod         - Run in production mode"
    Write-Host "  stop             - Stop running services"
    Write-Host "  docker-build     - Build Docker images"
    Write-Host "  docker-run       - Run Docker containers"
    Write-Host "  docker-stop      - Stop Docker containers"
    Write-Host "  docker-clean     - Clean Docker resources"
    Write-Host "  docker-dev       - Run Docker containers in development mode"
    Write-Host "  docker-prod      - Run Docker containers in production mode"
}

function Build {
    Write-Host "Building binaries..."
    go mod download
    go mod verify
    go build -o "$BUILD_DIR/agent.exe" cmd/agent/main.go
    go build -o "$BUILD_DIR/orchestrator.exe" cmd/orchestrator/main.go
}

function Clean {
    Write-Host "Cleaning up..."
    go clean
    if (Test-Path $BUILD_DIR) { Remove-Item -Recurse -Force $BUILD_DIR }
}

function Test {
    Write-Host "Running tests..."
    go test -v -cover ./...
}

function RunDev {
    Write-Host "Starting in development mode..."
    Build
    Start-Process "$BUILD_DIR/orchestrator.exe" -NoNewWindow
    Start-Process "$BUILD_DIR/agent.exe" -NoNewWindow
}

function RunProd {
    Write-Host "Starting in production mode..."
    Build
    Start-Process "$BUILD_DIR/orchestrator.exe" -NoNewWindow
    Start-Process "$BUILD_DIR/agent.exe" -NoNewWindow
}

function Stop {
    Write-Host "Stopping services..."
    Get-Process | Where-Object { $_.ProcessName -match 'agent|orchestrator' } | Stop-Process -Force
}

function Docker-Build {
    Write-Host "Building Docker images..."
    docker-compose build
}

function Docker-Run {
    Write-Host "Starting Docker containers..."
    docker-compose up -d
}

function Docker-Stop {
    Write-Host "Stopping Docker containers..."
    docker-compose down
}

function Docker-Clean {
    Write-Host "Cleaning Docker resources..."
    docker-compose down --rmi all --volumes --remove-orphans
}

function Docker-Dev {
    Write-Host "Starting Docker containers in development mode..."
    $env:COMPUTING_POWER = 2
    $env:TIME_ADDITION_MS = 100
    $env:TIME_SUBTRACTION_MS = 100
    $env:TIME_MULTIPLICATIONS_MS = 200
    $env:TIME_DIVISIONS_MS = 200
    docker-compose up -d
}

function Docker-Prod {
    Write-Host "Starting Docker containers in production mode..."
    $env:COMPUTING_POWER = 8
    docker-compose up -d
}

# Command Execution
switch ($Command) {
    "help" { Show-Help }
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "run-dev" { RunDev }
    "run-prod" { RunProd }
    "stop" { Stop }
    "docker-build" { Docker-Build }
    "docker-run" { Docker-Run }
    "docker-stop" { Docker-Stop }
    "docker-clean" { Docker-Clean }
    "docker-dev" { Docker-Dev }
    "docker-prod" { Docker-Prod }
    default { Write-Host "Unknown command. Use 'help' to see available commands." }
}