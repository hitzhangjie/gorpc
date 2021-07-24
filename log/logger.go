package log

// Logger defines logger behavior
type Logger interface {
	Trace(s string, v ...interface{})
	Debug(s string, v ...interface{})
	Info(s string, v ...interface{})
	Warn(s string, v ...interface{})
	Error(s string, v ...interface{})
	Fatal(s string, v ...interface{})
	WithPrefix(s string, v ...interface{}) Logger
}
