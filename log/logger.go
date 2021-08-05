package log

// Logger defines logger behavior
type Logger interface {
	Trace(s string, v ...interface{})
	Debug(s string, v ...interface{})
	Info(s string, v ...interface{})
	Warn(s string, v ...interface{})
	Error(s string, v ...interface{})
	Fatal(s string, v ...interface{})
	Flush() error
	WithPrefix(s string, v ...interface{}) Logger
}

// Level defines the levels of logging messages
type Level uint8

const (
	Trace Level = iota // print verbose debug message for framework
	Debug              // print verbose debug message for application
	Info               // print info message
	Warn               // print warn message
	Error              // print error message
	Fatal              // print fatal message and exit
)

func (l Level) String() string {
	switch l {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "Unknown"
	}
}

// RollType file rolling type
type RollType int

const (
	RollNONE = iota
	RollByFileSZ
	RollByDay
)
