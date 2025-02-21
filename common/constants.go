// Package common provides shared constants used across the application.
package common

// Error messages used throughout the application.
const (
	ErrInvalidRequestBody        = "Invalid request body"
	ErrExpressionNotFound        = "Expression not found"
	ErrTaskNotFound              = "Task not found"
	ErrFailedInitLogger          = "Failed to initialize logger: %v"
	ErrFailedSyncLogger          = "Failed to sync logger: %v"
	ErrFailedStartAgent          = "Failed to start agent"
	ErrFailedUpdateExpr          = "Failed to update expression error"
	ErrFailedCloseRespBody       = "Failed to close response body"
	ErrUnexpectedStatusCode      = "unexpected status code: %d"
	ErrFailedInitConfig          = "Failed to initialize config"
	ErrUnexpectedToken           = "unexpected token"
	ErrDivisionByZero            = "division by zero"
	ErrModuloByZero              = "modulo by zero"
	ErrInvalidModulo             = "modulo operation requires integer operands"
	ErrUnexpectedEndExpr         = "unexpected end of expression"
	ErrMissingCloseParen         = "missing closing parenthesis"
	ErrInvalidNumber             = "invalid number: %s"
	ErrExpressionNotFoundStorage = "expression not found"
	ErrFailedProcessExpression   = "Failed to process expression"
	ErrFailedProcessResult       = "Failed to process result"
	ErrFailedStartServer         = "Failed to start server"
	ErrServerShutdownFailed      = "Server shutdown failed"
	ErrEmptyExpressionID         = "expression ID cannot be empty"
	ErrInvalidStatusTransition   = "invalid status transition from %s to %s"
)

// Log messages used for logging application events.
const (
	LogProcessingExpression       = "Processing expression"
	LogTaskRetrieved              = "Task retrieved"
	LogExpressionRetrieved        = "Expression retrieved"
	LogAgentStarting              = "Starting agent"
	LogAgentStarted               = "Agent service started successfully"
	LogAgentStoppedGrace          = "Agent service stopped gracefully"
	LogAgentStopped               = "Agent stopped"
	LogWorkerStarting             = "Starting worker"
	LogWorkerStopped              = "Worker stopped"
	LogProcessingTask             = "Processing task"
	LogFailedProcessTask          = "Failed to process task"
	LogFailedGetTask              = "failed to get task"
	LogFailedSendResult           = "failed to send result"
	LogEmptyExpression            = "Empty expression received"
	LogFailedDecodeBody           = "Failed to decode request body"
	LogFailedSaveExpr             = "Failed to save expression"
	LogExpressionReceived         = "Expression received for calculation"
	LogNoTasksAvailable           = "No tasks available"
	LogFailedDecodeTask           = "Failed to decode task result"
	LogFailedUpdateTask           = "Failed to update task result"
	LogFailedGetTaskResult        = "Failed to get task after updating result"
	LogFailedUpdateExpr           = "Failed to update expression result"
	LogTaskProcessed              = "Task result processed successfully"
	LogOrchestratorStarted        = "Orchestrator service started successfully"
	LogOrchestratorStoppedGrace   = "Orchestrator service stopped gracefully"
	LogFailedSaveEmptyID          = "Failed to save expression: empty ID"
	LogExpressionSaved            = "Expression saved successfully"
	LogInvalidStatusTransition    = "Invalid status transition"
	LogExpressionStatusUpdated    = "Expression status updated"
	LogFailedUpdateStatusNotFound = "Failed to update expression status: expression not found"
	LogListedAllExpressions       = "Listed all expressions"
	LogExpressionNotFound         = "Expression not found"
	LogFailedParseExpression      = "Failed to parse expression"
	LogNoValidTasksCreated        = "No valid tasks created"
	LogFailedSaveTask             = "Failed to save task"
	LogTasksCreated               = "Tasks created successfully"
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

// Error format patterns used for wrapping errors.
const (
	ErrFormatWithWrap = "%s: %w"
)
