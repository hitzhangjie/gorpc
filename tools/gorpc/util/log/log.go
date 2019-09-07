package log

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	COLOR_RESET = "\033[0m"
	COLOR_GREEN = "\033[1;32m"
	COLOR_RED   = "\033[1;31m"
	COLOR_PINK  = "\033[1;35m"

	verbose bool
)

func InitLogging(v bool) {
	verbose = v
}

// Info print logging info at level INFO, if flag verbose true, filename and lineno will be logged.
func Info(format string, vals ...interface{}) {
	fn, _ := callerAddress(3)
	if verbose {
		fmt.Printf("%s[Info][%s] %s%s\n", COLOR_GREEN, fn, fmt.Sprintf(format, vals...), COLOR_RESET)
	} else {
		fmt.Printf("%s[Info] %s%s\n", COLOR_GREEN, fmt.Sprintf(format, vals...), COLOR_RESET)
	}
}

// Debug print logging info at level DEBUG, if flag verbose true, filename and lineno will be logged.
func Debug(format string, vals ...interface{}) {
	fn, _ := callerAddress(3)
	if verbose {
		fmt.Printf("%s[Debug][%s] %s%s\n", COLOR_PINK, fn, fmt.Sprintf(format, vals...), COLOR_RESET)
	}
}

// Error print logging info at level ERROR, if flag verbose true, filename and lineno will be logged.
func Error(format string, vals ...interface{}) {
	fn, _ := callerAddress(3)
	if verbose {
		fmt.Printf("%s[Error][%s] %s%s\n", COLOR_RED, fn, fmt.Sprintf(format, vals...), COLOR_RESET)
	} else {
		fmt.Printf("%s[Error] %s%s\n", COLOR_RED, fmt.Sprintf(format, vals...), COLOR_RESET)
	}
}

// callerAddress skip N level to get the caller's filename and lineno, if no caller return error.
func callerAddress(skip int) (string, error) {

	fpcs := make([]uintptr, 1)
	// Skip N levels to get the caller
	n := runtime.Callers(skip, fpcs)
	if n == 0 {
		return "", fmt.Errorf("MSG: NO CALLER")
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return "", fmt.Errorf("MSG: CALLER IS NIL")
	}

	// Print the file name and line number
	fileName, lineNo := caller.FileLine(fpcs[0] - 1)
	baseName := fileName[strings.LastIndex(fileName, "/")+1:]

	return fmt.Sprintf("%s:%d", baseName, lineNo), nil
}
