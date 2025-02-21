#!/usr/bin/env pwsh

param(
    [string]$Command = "run-dev"  # Default command argument
)

# Load variables from .env file
$envFile = ".env"
$DEFAULT_ENV = @{
    "ORCHESTRATOR_URL" = "http://localhost:8080"
    "PORT" = "8080"
}

if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            $DEFAULT_ENV[$key] = $value
        }
    }
}

# Set environment variables if not already set
foreach ($key in $DEFAULT_ENV.Keys) {
    if (-not (Test-Path env:$key)) {
        Set-Item -Path env:$key -Value $DEFAULT_ENV[$key]
    }
}

$BUILD_DIR = "build"  # Binary build directory
if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

# Function to display help information
function Show-Help {
    Write-Host "Available commands:"
    Write-Host "  all              - clean, deps, build"
    Write-Host "  deps             - download and verify dependencies"
    Write-Host "  check-deps       - tidy and verify modules"
    Write-Host "  build            - build binaries"
    Write-Host "  clean            - clean build artifacts"
    Write-Host "  test             - run tests with race and coverage"
    Write-Host "  test-short       - run short tests"
    Write-Host "  test-coverage    - generate coverage report"
    Write-Host "  lint             - run golangci-lint"
    Write-Host "  run              - run with .env variables"
    Write-Host "  run-dev          - run in development mode"
    Write-Host "  run-prod         - run in production mode"
    Write-Host "  stop             - stop running services"
    Write-Host "  run-agent        - run agent only"
    Write-Host "  run-orchestrator - run orchestrator only"
    Write-Host "  run-race         - run services with race detection"
}

# Function to build binaries
function Build {
    Write-Host "Building application..."
    go mod download
    go mod verify
    go build -o "$BUILD_DIR\agent.exe" cmd\agent\main.go
    go build -o "$BUILD_DIR\orchestrator.exe" cmd\orchestrator\main.go
}

# Function to clean builds
function Clean {
    Write-Host "Cleaning up..."
    go clean
    if (Test-Path $BUILD_DIR) {
        Remove-Item -Recurse -Force $BUILD_DIR
    }
}

# Function to manage dependencies
function Deps {
    Write-Host "Downloading dependencies..."
    go mod download
    go mod verify
}

# Function to check and tidy dependencies
function Check-Deps {
    Write-Host "Checking dependencies..."
    go mod tidy
    go mod verify
}

# Function to run tests
function Test {
    Write-Host "Running tests with race detection and coverage..."
    go test -v -race -cover ./...
}

# Function to run short tests
function Test-Short {
    Write-Host "Running short tests..."
    go test -v -short ./...
}

# Function to generate test coverage
function Test-Coverage {
    Write-Host "Generating test coverage report..."
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
}

# Function to run linter
function Lint {
    Write-Host "Running linter..."
    golangci-lint run ./...
}

# Function to run in development mode
function RunDev {
    Write-Host "Starting services in development mode..."
    
    $Env:COMPUTING_POWER = "2"
    $Env:TIME_ADDITION_MS = "100"
    $Env:TIME_SUBTRACTION_MS = "100"
    $Env:TIME_MULTIPLICATIONS_MS = "200"
    $Env:TIME_DIVISIONS_MS = "200"
    
    Build
    
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Get-ChildItem env: | Where-Object { $_.Name -match "COMPUTING_POWER|TIME_|ORCHESTRATOR_URL|PORT" }
}

# Function to run in production mode
function RunProd {
    Write-Host "Starting services in production mode..."
    
    $Env:COMPUTING_POWER = "8"
    $Env:TIME_ADDITION_MS = "1000"
    $Env:TIME_SUBTRACTION_MS = "1000"
    $Env:TIME_MULTIPLICATIONS_MS = "2000"
    $Env:TIME_DIVISIONS_MS = "2000"
    
    Build
    
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Get-ChildItem env: | Where-Object { $_.Name -match "COMPUTING_POWER|TIME_|ORCHESTRATOR_URL|PORT" }
}

# Function to run with .env variables
function Run {
    Write-Host "Starting services with .env variables..."
    Build
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
}

# Function to run agent only
function Run-Agent {
    Write-Host "Starting agent..."
    Build
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
}

# Function to run orchestrator only
function Run-Orchestrator {
    Write-Host "Starting orchestrator..."
    Build
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
}

# Function to run with race detection
function Run-Race {
    Write-Host "Starting services with race detection..."
    go build -race -o "$BUILD_DIR\agent.exe" cmd\agent\main.go
    go build -race -o "$BUILD_DIR\orchestrator.exe" cmd\orchestrator\main.go
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
}

# Function to stop services
function Stop {
    Write-Host "Stopping services..."
    Get-Process | Where-Object { $_.ProcessName -match 'agent|orchestrator' } | Stop-Process -Force
}

# Main execution switch
switch ($Command) {
    "help" { Show-Help }
    "all" { Clean; Deps; Build }
    "deps" { Deps }
    "check-deps" { Check-Deps }
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "test-short" { Test-Short }
    "test-coverage" { Test-Coverage }
    "lint" { Lint }
    "run" { Run }
    "run-dev" { RunDev }
    "run-prod" { RunProd }
    "stop" { Stop }
    "run-agent" { Run-Agent }
    "run-orchestrator" { Run-Orchestrator }
    "run-race" { Run-Race }
    default { Write-Host "Unknown command. Use 'help' to see available commands." }
}