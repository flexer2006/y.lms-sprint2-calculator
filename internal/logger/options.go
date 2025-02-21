package logger

// LogLevel определяет уровень логирования
type LogLevel string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
	Warn  LogLevel = "warn"
	Error LogLevel = "error"
)

// Options определяет настройки логгера
type Options struct {
	Level       LogLevel
	Encoding    string
	OutputPath  []string
	ErrorPath   []string
	Development bool
	LogDir      string
}

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
