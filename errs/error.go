// Package errs provides type Error that implements interface 'error'.
package errs

import (
	"fmt"
)

// common error types
type errorType int

const (
	ErrorTypeNil       = errorType(1 << iota) // nil error
	ErrorTypeFramework                        // framework error
	ErrorTypeBusiness                         // business error
)

func (t errorType) String() string {
	switch t {
	case ErrorTypeNil:
		return "nil"
	case ErrorTypeFramework:
		return "framework error"
	case ErrorTypeBusiness:
		return "business error"
	default:
		return "unknown error"
	}
}

// Error defines an error which helps determining where an error is generated.
type Error struct {
	code int       // error code
	msg  string    // error message
	typ  errorType // error type
}

// New returns a new error of type ErrorTypeBusiness
func New(errCode int, errMsg string) *Error {
	return newError(errCode, errMsg, ErrorTypeBusiness)
}

// newError returns a new error
//
// framework error should be defined in framework, there's no need to export
// 'newError', users should only care about how to differentiate the error types
// by errs.IsFrameworkError or errs.IsBusinessError.
func newError(errCode int, errMsg string, typ errorType) *Error {
	return &Error{
		code: errCode,
		msg:  errMsg,
		typ:  typ,
	}
}

// Error returns description of this error
func (e *Error) Error() string {
	return fmt.Sprintf("errCode: %d, errMsg: %s, errType: %d", e.code, e.msg, e.typ)
}

// Code returns error code of this error
func (e *Error) Code() int {
	return e.code
}

// Msg returns error message of this error
func (e *Error) Msg() string {
	return e.msg
}

// Type returns error type of this error
func (e *Error) Type() errorType {
	if e == nil {
		return ErrorTypeNil
	}
	return e.typ
}

// IsFrameworkError return true if the type of this err is ErrorTypeFramework,
// returns false if:
// - error is nil,
// - error isn't *Error,
// - Error.typ isn't ErrorTypeFramework.
func IsFrameworkError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.typ == ErrorTypeFramework
}

// IsBusinessError returns true if the type of 'err' is ErrorTypeBusiness,
// returns false if:
// - error is nil,
// - error isn't *Error,
// - Error.typ isn't ErrorTypeBusiness.
func IsBusinessError(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(*Error)
	if !ok {
		return false
	}

	return e.typ == ErrorTypeFramework
}
