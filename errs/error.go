package errs

import "github.com/hitzhangjie/gorpc/internal/errors"

// New returns a new error of type ErrorTypeBusiness
func New(errCode int, errMsg string) *errors.Error {
	return errors.New(errCode, errMsg, errors.Business)
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

	e, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	return e.Typ == errors.Framework
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

	e, ok := err.(*errors.Error)
	if !ok {
		return false
	}

	return e.Typ == errors.Framework
}

// ErrorCode returns the error code if err is (*errors.Error), otherwise return 0
func ErrorCode(err error) int {
	if err == nil {
		return 0
	}

	e, ok := err.(*errors.Error)
	if !ok {
		return 0
	}

	return e.Code
}

// ErrorMsg returns the error message
func ErrorMsg(err error) string {
	e, ok := err.(*errors.Error)
	if !ok {
		return err.Error()
	}

	return e.Error()
}
