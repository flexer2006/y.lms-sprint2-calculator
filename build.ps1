#!/usr/bin/env pwsh

param(
    [string]$Command = "run-dev"  # Аргумент командной строки, задающий действие (по умолчанию run-dev)
)

# Загружаем переменные из .env файла
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

# Устанавливаем переменные окружения, если они не заданы
foreach ($key in $DEFAULT_ENV.Keys) {
    if (-not (Test-Path env:$key)) {
        Set-Item -Path env:$key -Value $DEFAULT_ENV[$key]
    }
}

$BUILD_DIR = "build"  # Папка для сборки бинарных файлов
if (-not (Test-Path $BUILD_DIR)) {
    New-Item -ItemType Directory -Path $BUILD_DIR | Out-Null
}

# Функция для сборки бинарных файлов
function Build {
    Write-Host "Building application..."
    go mod download  # Загружаем зависимости
    go mod verify    # Проверяем целостность модулей
    go build -o "$BUILD_DIR\agent.exe" cmd\agent\main.go
    go build -o "$BUILD_DIR\orchestrator.exe" cmd\orchestrator\main.go
}

# Функция для очистки билдов
function Clean {
    Write-Host "Cleaning up..."
    go clean  # Очищаем кеш Go
    if (Test-Path $BUILD_DIR) {
        Remove-Item -Recurse -Force $BUILD_DIR
    }
}

# Функция запуска в режиме разработки
function RunDev {
    Write-Host "Starting services in development mode..."
    
    # Устанавливаем переменные окружения для dev-режима
    $Env:COMPUTING_POWER = "2"
    $Env:TIME_ADDITION_MS = "100"
    $Env:TIME_SUBTRACTION_MS = "100"
    $Env:TIME_MULTIPLICATIONS_MS = "200"
    $Env:TIME_DIVISIONS_MS = "200"
    
    Build  # Сборка приложения
    
    # Запускаем процессы
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Get-ChildItem env: | Where-Object { $_.Name -match "COMPUTING_POWER|TIME_|ORCHESTRATOR_URL|PORT" }
}

# Функция запуска в продакшн-режиме
function RunProd {
    Write-Host "Starting services in production mode..."
    
    # Устанавливаем переменные окружения для продакшн-режима
    $Env:COMPUTING_POWER = "8"
    $Env:TIME_ADDITION_MS = "1000"
    $Env:TIME_SUBTRACTION_MS = "1000"
    $Env:TIME_MULTIPLICATIONS_MS = "2000"
    $Env:TIME_DIVISIONS_MS = "2000"
    
    Build  # Сборка приложения
    
    # Запускаем процессы
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables set to:"
    Get-ChildItem env: | Where-Object { $_.Name -match "COMPUTING_POWER|TIME_|ORCHESTRATOR_URL|PORT" }
}

# Новая функция запуска с переменными из .env
function Run {
    Write-Host "Starting services with .env variables..."
    
    Build  # Сборка приложения
    
    # Запускаем процессы
    Start-Process -FilePath "$BUILD_DIR\orchestrator.exe" -NoNewWindow
    Start-Process -FilePath "$BUILD_DIR\agent.exe" -NoNewWindow
    
    Write-Host "Environment variables loaded from .env (or defaults):"
    Get-ChildItem env: | Where-Object { $_.Name -match "COMPUTING_POWER|TIME_|ORCHESTRATOR_URL|PORT" }
}

# Функция остановки запущенных процессов
function Stop {
    Write-Host "Stopping services..."
    Get-Process -Name "agent", "orchestrator" -ErrorAction SilentlyContinue | Stop-Process
}

# Функция запуска тестов
function Test {
    Write-Host "Running tests..."
    go test -v -race -cover ./...
}

# Обрабатываем переданный аргумент и вызываем соответствующую функцию
switch ($Command) {
    "build" { Build }
    "clean" { Clean }
    "test" { Test }
    "run-dev" { RunDev }
    "run-prod" { RunProd }
    "run" { Run }  # Добавлена новая команда
    "stop" { Stop }
    default { Write-Host "Unknown command: $Command" }
}