package logger

// LogLevel defines the logging level.
type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

// Options defines the logger settings.
type Options struct {
	Level       LogLevel
	Encoding    string
	OutputPath  []string
	ErrorPath   []string
	Development bool
	LogDir      string
}

// DefaultOptions returns the default logger options.
func DefaultOptions() Options {
	return Options{
		Level:       Info,
		Encoding:    "json",
		OutputPath:  []string{"stdout"},
		ErrorPath:   []string{"stderr"},
		Development: false,
		LogDir:      "logs",
	}
}
