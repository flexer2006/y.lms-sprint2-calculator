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

# Command Execution
switch ($Command) {
    "help" { Show-Help }
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "run-dev" { RunDev }
    "run-prod" { RunProd }
    "stop" { Stop }
    default { Write-Host "Unknown command. Use 'help' to see available commands." }
}