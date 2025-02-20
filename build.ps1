#!/usr/bin/env pwsh

param(
    [string]$Command = "run-dev"
)

# Default environment variables
$DEFAULT_ENV = @{
    "COMPUTING_POWER" = "4"
    "TIME_ADDITION_MS" = "1000"
    "TIME_SUBTRACTION_MS" = "1000"
    "TIME_MULTIPLICATIONS_MS" = "2000"
    "TIME_DIVISIONS_MS" = "2000"
    "ORCHESTRATOR_URL" = "http://localhost:8080"
    "PORT" = "8080"
}

# Set default environment variables if not already set
foreach ($key in $DEFAULT_ENV.Keys) {
    if (-not (Test-Path env:$key)) {
        Set-Item -Path env:$key -Value $DEFAULT_ENV[$key]
    }
}

$BUILD_DIR = "build"
if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

function Build {
    Write-Host "Building application..."
    go mod download
    go mod verify
    go build -o "$BUILD_DIR\agent.exe" cmd\agent\main.go
    go build -o "$BUILD_DIR\orchestrator.exe" cmd\orchestrator\main.go
}

function Clean {
    Write-Host "Cleaning up..."
    go clean
    if (Test-Path $BUILD_DIR) {
        Remove-Item -Recurse -Force $BUILD_DIR
    }
}

function RunDev {
    Write-Host "Starting services in development mode..."
    # Development environment variables
    $Env:COMPUTING_POWER = "2"
    $Env:TIME_ADDITION_MS = "100"
    $Env:TIME_SUBTRACTION_MS = "100"
    $Env:TIME_MULTIPLICATIONS_MS = "200"
    $Env:TIME_DIVISIONS_MS = "200"
    
    Build
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Write-Host "COMPUTING_POWER: $($Env:COMPUTING_POWER)"
    Write-Host "TIME_ADDITION_MS: $($Env:TIME_ADDITION_MS)"
    Write-Host "TIME_SUBTRACTION_MS: $($Env:TIME_SUBTRACTION_MS)"
    Write-Host "TIME_MULTIPLICATIONS_MS: $($Env:TIME_MULTIPLICATIONS_MS)"
    Write-Host "TIME_DIVISIONS_MS: $($Env:TIME_DIVISIONS_MS)"
    Write-Host "ORCHESTRATOR_URL: $($Env:ORCHESTRATOR_URL)"
    Write-Host "PORT: $($Env:PORT)"
}

function RunProd {
    Write-Host "Starting services in production mode..."
    # Production environment variables
    $Env:COMPUTING_POWER = "8"
    $Env:TIME_ADDITION_MS = "1000"
    $Env:TIME_SUBTRACTION_MS = "1000"
    $Env:TIME_MULTIPLICATIONS_MS = "2000"
    $Env:TIME_DIVISIONS_MS = "2000"
    
    Build
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Write-Host "COMPUTING_POWER: $($Env:COMPUTING_POWER)"
    Write-Host "TIME_ADDITION_MS: $($Env:TIME_ADDITION_MS)"
    Write-Host "TIME_SUBTRACTION_MS: $($Env:TIME_SUBTRACTION_MS)"
    Write-Host "TIME_MULTIPLICATIONS_MS: $($Env:TIME_MULTIPLICATIONS_MS)"
    Write-Host "TIME_DIVISIONS_MS: $($Env:TIME_DIVISIONS_MS)"
    Write-Host "ORCHESTRATOR_URL: $($Env:ORCHESTRATOR_URL)"
    Write-Host "PORT: $($Env:PORT)"
}

function Stop {
    Write-Host "Stopping services..."
    Get-Process -Name "agent", "orchestrator" -ErrorAction SilentlyContinue | Stop-Process
}

function Test {
    Write-Host "Running tests..."
    go test -v -race -cover ./...
}

switch ($Command) {
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "run-dev" { RunDev }
    "run-prod" { RunProd }
    "stop" { Stop }
    default { Write-Host "Unknown command: $Command" }
} 