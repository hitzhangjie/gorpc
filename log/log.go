package log

type LogLevel uint8

const (
	// Trace trace level logging, usually for debugging framework or components internal events
	LogLevelTrace LogLevel = iota

	// LogLevelDebug debug level logging
	LogLevelDebug

	// LogLevelInfo info level logging
	LogLevelInfo

	// LogLevelWarn warn level logging
	LogLevelWarn

	// Error error level logging
	LogLevelError

	// LogLevelFatal fatal level logging, please pay attention fatal error generates a panic
	LogLevelFatal
)

// Logger defines logging methods and behaviors
type Logger interface {
	// Trace show trace and higher level logging info
	Trace(fmt string, args ...interface{})

	// Debug show debug and higher level logging info
	Debug(fmt string, args ...interface{})

	// Info show info and higher level logging info
	Info(fmt string, args ...interface{})

	// Warn show warn and higher level logging info
	Warn(fmt string, args ...interface{})

	// Error show error and higher level logging info
	Error(fmt string, args ...interface{})

	// Fatal show fatal level logging info and panic
	Fatal(fmt string, args ...interface{})
}
