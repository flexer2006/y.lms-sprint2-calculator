# Distributed Arithmetic Expression Calculator

## Description
This is a distributed calculator project that processes arithmetic expressions. The system consists of two main components: an orchestrator and agents. The orchestrator service manages incoming calculation requests and distributes the workload among multiple agent instances. Agents perform the actual arithmetic calculations with configurable processing times for different operations (addition, subtraction, multiplication, division).

## Features

Based on test files and implementation, this calculator has the following capabilities:

1. Arithmetic Operations: Basic operations: addition (+), subtraction (-), multiplication (*), division (/), high-precision decimal number processing, support for very large and very small numbers, proper operator precedence handling.
2. Expression Functions: Support for parentheses in nested expressions: (2 + 3) * (4 + 5), unary minus operator in different contexts (-2, 2 * -3), multiple operations in one expression, complex nested expressions, flexible whitespace handling.
3. Input Validation: Empty expression validation, balanced parentheses checking, decimal point usage validation, invalid character prevention, consecutive operator checking, missing operand/operator validation, division by zero protection.
4. Distributed Processing: Parallel computation processing, task distribution across multiple agents, configurable operation timing, request/response logging, error handling and status tracking.
5. Additional Features: Expression status tracking (waiting, in progress, completed, error), detailed error reporting, comprehensive logging system, long expression support, high-precision decimal calculations.

## Technologies Used

- **Go v1.23.6**
- **Go Standard Libraries**
- **Zap logger v1.27.0**
- **Google UUID Generator v1.6.0**
- **Testify for Enhanced Testing v1.10.0**
- **HTTP Router gorilla/mux v1.8.1**
- **Debug Output go-spew v1.1.1**
- **String Difference Calculation and Formatting - go-difflib v1.0.0**
- **Multiple Error Combining Library - multierr v1.11.0**
- **YAML Processing Library - yaml.v3 v3.0.1**

## Installation

### Clone the Repository

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

**Build the images:**
```
docker build -t calculator-orchestrator -f Dockerfile --target orchestrator .
```
```
docker build -t calculator-agent -f Dockerfile --target agent .
```

**Run the orchestrator:**
```
docker run -d -p 8080:8080 --name orchestrator calculator-orchestrator
```
**Run one or more agents:**
```
docker run -d --name agent1 calculator-agent
```

## Method 2: Using Makefile

**View available commands:**
```
make help
```
**Build both services:**
```
make build
```
**Run in any convenient mode:**
```
make run
```
```
make run-dev
```
```
make run-prod
```
```
make run-agent
```
```
make run-orchestrator
```
```
make run-race
```

## Method 3: Using PowerShell

**Build the service with PowerShell:**
```
.\build.ps1
```
**Run in any convenient mode:**
```
run
```
```
run-dev
```
```
run-prod
```

## Method 4: Using CMD

**Build the project:**
```
go build -o bin/orchestrator.exe cmd/orchestrator/main.go
```
```
go build -o bin/agent.exe cmd/agent/main.go
```
**Run the services:**
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
**Check calculation result:**
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

### List All Expressions

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

### Testing

**To run tests:**
```
go test ./...
```
**For detailed test output:**
```
go test ./... -v
```