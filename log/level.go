package log

const (
	Trace Level = iota // print verbose debug message for framework
	Debug              // print verbose debug message for application
	Info               // print info message
	Warn               // print warn message
	Error              // print error message
	Fatal              // print fatal message and exit
)

// Level defines the levels of logging messages
type Level uint8

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
