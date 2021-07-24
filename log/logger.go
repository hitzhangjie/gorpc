package log

// Level defines the levels of logging messages
type Level uint8

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

// Logger defines logger behavior
type Logger interface {
	Trace(fmt string, args ...interface{})
	Debug(fmt string, args ...interface{})
	Info(fmt string, args ...interface{})
	Warn(fmt string, args ...interface{})
	Error(fmt string, args ...interface{})
	Fatal(fmt string, args ...interface{})
	WithPrefix(fmt string, v...interface{})
}
