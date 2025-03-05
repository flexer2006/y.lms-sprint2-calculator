// Package common provides shared constants used across the application.
package common

// Error messages used throughout the application.
const (
	ErrInvalidRequestBody      = "Invalid request body"
	ErrExpressionNotFound      = "Expression not found"
	ErrTaskNotFound            = "Task not found"
	ErrFailedInitLogger        = "Failed to initialize logger: %v"
	ErrFailedSyncLogger        = "Failed to sync logger: %v"
	ErrFailedStartAgent        = "Failed to start agent"
	ErrFailedCloseRespBody     = "Failed to close response body"
	ErrUnexpectedStatusCode    = "unexpected status code: %d"
	ErrFailedInitConfig        = "Failed to initialize config"
	ErrUnexpectedToken         = "unexpected token"
	ErrDivisionByZero          = "division by zero"
	ErrModuloByZero            = "modulo by zero"
	ErrInvalidModulo           = "modulo operation requires integer operands"
	ErrUnexpectedEndExpr       = "unexpected end of expression"
	ErrMissingCloseParen       = "missing closing parenthesis"
	ErrFailedProcessExpression = "Failed to process expression"
	ErrFailedProcessResult     = "Failed to process result"
	ErrFailedStartServer       = "Failed to start server"
	ErrServerShutdownFailed    = "Server shutdown failed"
)

// Log messages used for logging application events.
const (
	LogTaskRetrieved              = "Task retrieved"
	LogExpressionRetrieved        = "Expression retrieved"
	LogAgentStarted               = "Agent service started successfully"
	LogAgentStoppedGrace          = "Agent service stopped gracefully"
	LogFailedSendResult           = "failed to send result"
	LogNoTasksAvailable           = "No tasks available"
	LogFailedDecodeTask           = "Failed to decode task result"
	LogFailedUpdateTask           = "Failed to update task result"
	LogFailedGetTaskResult        = "Failed to get task after updating result"
	LogFailedUpdateExpr           = "Failed to update expression result"
	LogTaskProcessed              = "Task result processed successfully"
	LogOrchestratorStarted        = "Orchestrator service started successfully"
	LogOrchestratorStoppedGrace   = "Orchestrator service stopped gracefully"
	LogInvalidStatusTransition    = "Invalid status transition"
	LogExpressionStatusUpdated    = "Expression status updated"
	LogFailedUpdateStatusNotFound = "Failed to update expression status: expression not found"
	LogListedAllExpressions       = "Listed all expressions"
	LogFailedParseExpression      = "Failed to parse expression"
)

// HTTP headers and content types used in the application.
const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

// URL paths used for API endpoints.
const (
	PathTask         = "/task"
	PathInternalTask = "%s/internal/task"
)

// Field names used in JSON and other data structures.
const (
	FieldCount           = "count"
	FieldStatus          = "status"
	FieldExpressionID    = "expressionID"
	FieldOperation       = "operation"
	FieldTaskID          = "taskID"
	FieldNewStatus       = "newStatus"
	FieldOldStatus       = "oldStatus"
	FieldToken           = "token"
	FieldWorkerID        = "worker_id"
	FieldResult          = "result"
	FieldExpression      = "expression"
	FieldTraceID         = "trace_id"
	FieldCorrelationID   = "correlation_id"
	FieldTokens          = "tokens"
	FieldPosition        = "position"
	FieldRequestID       = "request_id"
	FieldPort            = "port"
	FieldComputingPower  = "computing_power"
	FieldOrchestratorURL = "orchestrator_url"
	FieldID              = "id"
)

// Parser log messages used during expression parsing.
const (
	LogUnexpectedEndExpr      = "Unexpected end of expression"
	LogFailedParseParentheses = "Failed to parse expression in parentheses"
	LogMissingCloseParen      = "Missing closing parenthesis"
	LogFailedParseNegative    = "Failed to parse negative factor"
	LogInvalidNumberFormat    = "Invalid number format"
	LogUnexpectedToken        = "Unexpected token"
)

// Logger field names used in structured logging.
const (
	LogFieldTimestamp  = "timestamp"
	LogFieldLevel      = "level"
	LogFieldLogger     = "logger"
	LogFieldCaller     = "caller"
	LogFieldMessage    = "message"
	LogFieldStacktrace = "stacktrace"
)

const (
	ErrFormatWithWrap = "%s: %w"
)
