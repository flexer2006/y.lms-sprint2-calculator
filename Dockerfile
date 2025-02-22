# --- Stage: Builder ---
    FROM golang:1.23-alpine AS builder

    WORKDIR /app
    
    # Copy dependency files and download them
    COPY go.mod go.sum ./
    RUN go mod download
    
    # Copy the source code
    COPY . .
    
    # Build Linux-AMD64 binaries for agent and orchestrator
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/agent-linux-amd64 cmd/agent/main.go && \
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/orchestrator-linux-amd64 cmd/orchestrator/main.go
    
    # --- Final Stage: orchestrator ---
    FROM alpine:latest AS orchestrator
    
    RUN apk --no-cache add ca-certificates && \
        mkdir -p /app/logs/orchestrator
    
    WORKDIR /app
    
    # Copy the orchestrator binary from builder
    COPY --from=builder /app/build/orchestrator-linux-amd64 /app/orchestrator
    
    # Set environment variables for orchestrator
    ENV PORT=8080
    
    # Start orchestrator
    CMD ["/app/orchestrator"]
    
    # --- Final Stage: agent ---
    FROM alpine:latest AS agent
    
    RUN apk --no-cache add ca-certificates && \
        mkdir -p /app/logs/agent
    
    WORKDIR /app
    
    # Copy the agent binary from builder
    COPY --from=builder /app/build/agent-linux-amd64 /app/agent
    
    # Set environment variables for agent (values are taken from your .env file)
    ENV COMPUTING_POWER=4 \
        TIME_ADDITION_MS=1000 \
        TIME_SUBTRACTION_MS=1000 \
        TIME_MULTIPLICATIONS_MS=2000 \
        TIME_DIVISIONS_MS=2000 \
        ORCHESTRATOR_URL=http://localhost:8080
    
    # Start agent
    CMD ["/app/agent"]
    