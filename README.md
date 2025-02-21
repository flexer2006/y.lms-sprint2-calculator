# Calculator API

This service provides a REST API for calculating arithmetic expressions asynchronously.

## API Examples

### Calculate Expression

#### Successful Calculation

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "2+2*2"
}'
```

Response (HTTP 201 Created):
```json
{
    "id": "123e4567-e89b-12d3-a456-426614174000"
}
```

#### Check Calculation Result

```bash
curl --location 'http://localhost:8080/api/v1/expressions/123e4567-e89b-12d3-a456-426614174000'
```

Response (HTTP 200 OK):
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

### Error Cases

#### Empty Expression

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": ""
}'
```

Response (HTTP 422 Unprocessable Entity):
```json
{
    "error": "Expression cannot be empty"
}
```

#### Invalid Expression Format

```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression": "2++2"
}'
```

Response (HTTP 201 Created, then check status):
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

#### Expression Not Found

```bash
curl --location 'http://localhost:8080/api/v1/expressions/non-existent-id'
```

Response (HTTP 404 Not Found):
```json
{
    "error": "Expression not found"
}
```

### List All Expressions

```bash
curl --location 'http://localhost:8080/api/v1/expressions'
```

Response (HTTP 200 OK):
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