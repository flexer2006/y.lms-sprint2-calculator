# --- Stage: Builder ---
    FROM golang:1.23-alpine AS builder

    WORKDIR /app
    
    # Копируем файлы зависимостей и загружаем их
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Копируем исходный код
    COPY . .
    
    # Сборка Linux-AMD64 бинарников для agent и orchestrator
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/agent-linux-amd64 cmd/agent/main.go && \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/orchestrator-linux-amd64 cmd/orchestrator/main.go
    
    # --- Final Stage: orchestrator ---
    FROM alpine:latest AS orchestrator
    
    RUN apk --no-cache add ca-certificates && \
        mkdir -p /app/logs/orchestrator
    
    WORKDIR /app
    
    # Копируем бинарник orchestrator из builder
    COPY --from=builder /app/build/orchestrator-linux-amd64 /app/orchestrator
    
    # Задаём переменные окружения для orchestrator
    ENV PORT=8080
    
    # Запуск orchestrator
    CMD ["/app/orchestrator"]
    
    # --- Final Stage: agent ---
    FROM alpine:latest AS agent
    
    RUN apk --no-cache add ca-certificates && \
        mkdir -p /app/logs/agent
    
    WORKDIR /app
    
    # Копируем бинарник agent из builder
    COPY --from=builder /app/build/agent-linux-amd64 /app/agent
    
    # Задаём переменные окружения для agent (значения берутся из вашего .env)
    ENV COMPUTING_POWER=4 \
        TIME_ADDITION_MS=1000 \
        TIME_SUBTRACTION_MS=1000 \
        TIME_MULTIPLICATIONS_MS=2000 \
        TIME_DIVISIONS_MS=2000 \
        ORCHESTRATOR_URL=http://localhost:8080
    
    # Запуск agent
    CMD ["/app/agent"]
    