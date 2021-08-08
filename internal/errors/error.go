// Package errors provides type Error that implements interface 'error'.
package errors

import (
	"fmt"
)

// common error types
type errorType int

const (
	errorTypeNil = errorType(iota) // nil error
	Framework                      // framework error
	Business                       // business error
)

func (t errorType) String() string {
	switch t {
	case errorTypeNil:
		return "nil"
	case Framework:
		return "framework error"
	case Business:
		return "business error"
	default:
		return "unknown error"
	}
}

// New returns a new error
//
// framework error should be defined in framework, there's no need to export
// 'newError', users should only care about how to differentiate the error types
// by errors.IsFrameworkError or errors.IsBusinessError.
func New(errCode int, errMsg string, typ errorType) *Error {
	return &Error{
		Code: errCode,
		Msg:  errMsg,
		Typ:  typ,
	}
}

// Error defines an error which helps determining where an error is generated.
type Error struct {
	Code int       // error code
	Msg  string    // error message
	Typ  errorType // error type
}

// Error returns description of this error
func (e *Error) Error() string {
	return fmt.Sprintf("errCode: %d, errMsg: %s, errType: %s", e.Code, e.Msg, e.Typ)
}
