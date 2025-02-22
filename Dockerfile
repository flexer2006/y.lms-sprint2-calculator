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

# Define ARGs for OS and architecture (they будут передаваться через build аргументы)
ARG TARGETOS
ARG TARGETARCH

# Build for the specified OS and architecture
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/build/agent-$TARGETOS-$TARGETARCH cmd/agent/main.go
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/build/orchestrator-$TARGETOS-$TARGETARCH cmd/orchestrator/main.go

# Final stage for orchestrator
FROM alpine:latest AS orchestrator

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy orchestrator binary from builder, dynamically selecting based on target OS and architecture
COPY --from=builder /app/build/orchestrator-${TARGETOS}-${TARGETARCH} /app/orchestrator

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

# Copy agent binary from builder, dynamically selecting based on target OS and architecture
COPY --from=builder /app/build/agent-${TARGETOS}-${TARGETARCH} /app/agent

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
