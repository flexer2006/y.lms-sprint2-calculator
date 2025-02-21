# Distributed Arithmetic Expression Calculator

## Description
This is a distributed calculator project that processes arithmetic expressions. The system consists of two main components: orchestrator and agents. The orchestrator service manages incoming calculation requests and distributes the workload among multiple agent instances. Agents perform the actual arithmetic calculations with configurable processing times for different operations (addition, subtraction, multiplication, division).

## Features

Based on test files and implementation, this calculator has the following capabilities:

1. Arithmetic Operations: Basic operations: addition ( + ), subtraction ( - ), multiplication ( * ), division ( / ), high-precision decimal number processing, support for very large and very small numbers, correct operator precedence handling.
2. Expression Functions: Support for brackets in nested expressions: (2 + 3) * (4 + 5), unary minus operator in different contexts (-2, 2 * -3), multiple operations in one expression, complex nested expressions, flexible whitespace handling.
3. Input Validation: Empty expression validation, balanced bracket checking, decimal point usage validation, invalid character prevention, consecutive operator checking, missing operand/operator validation, division by zero protection.
4. Distributed Processing: Parallel computation processing, task distribution across multiple agents, operation time configuration for different operations, request/response logging, error handling and status tracking.
5. Additional Features: Expression status tracking (waiting, in progress, completed, error), detailed error reporting, comprehensive logging system, long expression support, high-precision decimal calculations.

## Prerequisites

Before installation, ensure you have:

1. **Go version 1.23.6 or higher**
   - Download Go from [official website](https://golang.org/dl/)
   - Verify installation: `go version`

2. **Git**
   - Install Git from [git-scm.com](https://git-scm.com/)
   - Verify installation: `git --version`

3. **Docker** (optional, for container deployment)
   - Install Docker from [docker.com](https://www.docker.com/)
   - Verify installation: `docker --version`

4. **Make** (optional, for using Makefile)
   - Windows: install via [Chocolatey](https://chocolatey.org/): `choco install make`
   - Verify installation: `make --version`

## Technologies Used

- **Go v1.23.6** - Main programming language
- **Go Standard Libraries** - Basic functionality
- **Zap logger v1.27.0** - High-performance logging
- **Google UUID generator v1.6.0** - Unique identifier generation
- **Testify for enhanced testing v1.10.0** - Improved testing
- **HTTP router gorilla/mux v1.8.1** - HTTP request routing
- **Debug output go-spew v1.1.1** - Enhanced debugging
- **For calculating and formatting string differences - go-difflib v1.0.0** - String comparison
- **Library for combining multiple errors into one - multierr v1.11.0** - Error handling
- **Library for working with YAML - yaml.v3 v3.0.1** - Configuration handling

## Installation

### Clone the repository

**Using HTTPS:**
```
git clone https://github.com/flexer2006/y.lms-sprint2-calculator.git
```

**Or using SSH:**
```
git clone git@github.com:flexer2006/y.lms-sprint2-calculator.git
```

## Running

### Navigate to the project directory:
```
cd y.lms-sprint2-calculator
```

### Method 1: Using Docker

#### 1. Build Images

**Build orchestrator image:**
```
docker build -t calculator-orchestrator -f Dockerfile --target orchestrator .
```
This command creates an image for the orchestrator service that manages calculation distribution.

**Build agent image:**
```
docker build -t calculator-agent -f Dockerfile --target agent .
```
This command creates an image for the agent service that performs calculations.

#### 2. Run Containers

**Start the orchestrator:**
```
docker run -d -p 8080:8080 --name orchestrator calculator-orchestrator
```
Flags:
- `-d`: run in background mode
- `-p 8080:8080`: port forwarding from container to host
- `--name orchestrator`: container name

**Start one or more agents:**
```
docker run -d --name agent1 calculator-agent
docker run -d --name agent2 calculator-agent
docker run -d --name agent3 calculator-agent
```
You can start as many agents as needed to handle the workload.

#### 3. Container Management

**View running containers:**
```
docker ps
```

**Stop containers:**
```
docker stop orchestrator agent1 agent2 agent3
```

**Remove containers:**
```
docker rm orchestrator agent1 agent2 agent3
```

### Method 2: Using Makefile

Makefile provides convenient commands for building and running the project.

#### Available commands:

1. **View all available commands:**
```
make help
```

2. **Build the project:**
```
make build
```
Builds both services (orchestrator and agent)

3. **Run in different modes:**

- **Standard run:**
```
make run
```
Starts orchestrator and one agent with default settings

- **Development mode run:**
```
make run-dev
```
Starts with extended logging and additional debug information

- **Production mode run:**
```
make run-prod
```
Starts with optimized settings for production environment

- **Run agent only:**
```
make run-agent
```

- **Run orchestrator only:**
```
make run-orchestrator
```

- **Run with race detection:**
```
make run-race
```
Starts with data race detector enabled

### Method 3: Using PowerShell

PowerShell scripts provide an easy way to build and run the project on Windows.

#### 1. Build the project

**Run the build script:**
```
.\build.ps1
```
This script will:
- Check dependencies
- Compile orchestrator and agent
- Prepare configuration files

#### 2. Run in different modes

**Standard run:**
```
.\run.ps1
```
Starts orchestrator and agent with default settings

**Development mode run:**
```
.\run-dev.ps1
```
Starts with extended logging

**Production mode run:**
```
.\run-prod.ps1
```
Starts with optimized settings

## Method 4: Using CMD

**Build the project:**
```
go build -o bin/orchestrator.exe cmd/orchestrator/main.go
```
```
go build -o bin/agent.exe cmd/agent/main.go
```
**Start the services:**
```
bin\orchestrator.exe
```
```
bin\agent.exe
```

## Environment Variables

### Orchestrator Service:

1. `PORT` - HTTP server port (default: 8080)
2. `TIME_ADDITION_MS` - Addition operation execution time (default: 100)
3. `TIME_SUBTRACTION_MS` - Subtraction operation execution time (default: 100)
4. `TIME_MULTIPLY_MS` - Multiplication operation execution time (default: 200)
5. `TIME_DIVISION_MS` - Division operation execution time (default: 200)

### Agent Service:

1. `COMPUTING_POWER` - Number of simultaneous calculations (default: 1)
2. `ORCHESTRATOR_URL` - Orchestrator service URL (default: http://localhost:8080)

## API Endpoints

### Orchestrator Service:

- `POST /api/v1/calculate` - Submit expression for calculation
- `GET /api/v1/expressions` - List all expressions
- `GET /api/v1/expressions/{id}` - Get expression status and result

### Internal API

- `GET /internal/task` - Get next task (used by agents)
- `POST /internal/task` - Submit task result (used by agents)

## Usage Examples

### Successful Calculation

**Using curl request:**
```
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":"2+2*2"}'
```
**Response (HTTP 201 Created):**
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000"
}
```
**Check example solution:**
```
curl --location 'http://localhost:8080/api/v1/expressions/123e4567-e89b-12d3-a456-426614174000'
```
**Response (HTTP 200 OK):**
```json
{
    "expression": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "expression": "2+2*2",
        "status": "COMPLETE",
        "result": 6
    }
}
```

### Unsuccessful Calculation

**Empty expression:**
```
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":""}'
```
**Response (HTTP 422 Unprocessable Entity):**
```json
{
    "error": "Expression cannot be empty"
}
```

**Invalid expression format:**
```
curl -L 'http://localhost:8080/api/v1/calculate' -H 'Content-Type: application/json' --data '{"expression":"2++2"}'
```
**Response (HTTP 201 Created, but check status):**
```json
{
    "expression": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "expression": "2++2",
        "status": "ERROR",
        "error": "invalid expression: too few tokens"
    }
}
```

**Expression not found:**
```
curl --location 'http://localhost:8080/api/v1/expressions/non-existent-id'
```
Response (HTTP 404 Not Found):
```json
{
    "error": "Expression not found"
}
```

### List all expressions

**Input:**
```
curl --location 'http://localhost:8080/api/v1/expressions'
```
**Response (HTTP 200 OK):**
```json
{
    "expressions": [
        {
            "id": "123e4567-e89b-12d3-a456-426614174000",
            "expression": "2+2*2",
            "status": "COMPLETE",
            "result": 6
        },
        {
            "id": "987fcdeb-51d3-12a4-b678-426614174000",
            "expression": "10-5",
            "status": "COMPLETE",
            "result": 5
        }
    ]
}
```

### Postman

You can use Postman to test the project:
![[Pasted image 20250221114952.png]]

## Development and Testing

### Project Structure

The project is organized as follows:
- `cmd/` - Entry points for orchestrator and agent
- `internal/` - Internal application logic
- `pkg/` - Reusable packages
- `tests/` - Test files
- `configs/` - Configuration files

### Testing

**Run all tests:**
```
go test ./...
```

**Detailed test output:**
```
go test ./... -v
```

**Run tests with race detection:**
```
go test -race ./...
```

**Run specific test:**
```
go test ./tests -run TestCalculator
```

**Run tests with coverage:**
```
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Troubleshooting

### Common Issues

1. **"connection refused" error**
   - Ensure orchestrator is running and accessible
   - Check URL correctness in agent configuration

2. **Compilation error**
   - Run `go mod tidy` to update dependencies
   - Check Go version (1.23.6 or higher required)

3. **Agent not connecting to orchestrator**
   - Check network settings
   - Ensure ports are not blocked by firewall