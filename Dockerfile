# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build for Linux (static binaries)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/agent-linux-amd64 cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/build/orchestrator-linux-amd64 cmd/orchestrator/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/build/agent-linux-arm64 cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o /app/build/orchestrator-linux-arm64 cmd/orchestrator/main.go

# Build for Windows (static binaries)
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /app/build/agent-windows-amd64.exe cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /app/build/orchestrator-windows-amd64.exe cmd/orchestrator/main.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o /app/build/agent-windows-arm64.exe cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o /app/build/orchestrator-windows-arm64.exe cmd/orchestrator/main.go

# Build for macOS (static binaries)
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o /app/build/agent-macos-amd64 cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o /app/build/orchestrator-macos-amd64 cmd/orchestrator/main.go
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o /app/build/agent-macos-arm64 cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o /app/build/orchestrator-macos-arm64 cmd/orchestrator/main.go

# Final stage for orchestrator
FROM alpine:latest AS orchestrator

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy only Linux amd64 orchestrator binary
COPY --from=builder /app/build/orchestrator-linux-amd64 /app/orchestrator

# Copy environment variables file
COPY .env ./

# Set default environment variables
ENV PORT=8080

# Expose the port
EXPOSE 8080

# Create directory for logs
RUN mkdir -p /app/logs/orchestrator

# Start orchestrator service
CMD ["/bin/sh", "-c", "source .env 2>/dev/null || true; /app/orchestrator"]

# Final stage for agent
FROM alpine:latest AS agent

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy only Linux amd64 agent binary
COPY --from=builder /app/build/agent-linux-amd64 /app/agent

# Copy environment variables file
COPY .env ./

# Set default environment variables
ENV COMPUTING_POWER=4 \
    TIME_ADDITION_MS=1000 \
    TIME_SUBTRACTION_MS=1000 \
    TIME_MULTIPLICATIONS_MS=2000 \
    TIME_DIVISIONS_MS=2000 \
    ORCHESTRATOR_URL=http://localhost:8080

# Create directory for logs
RUN mkdir -p /app/logs/agent

# Start agent service
CMD ["/bin/sh", "-c", "source .env 2>/dev/null || true; /app/agent"]