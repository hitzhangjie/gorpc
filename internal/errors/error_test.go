package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hitzhangjie/gorpc/internal/errors"
)

func TestNew(t *testing.T) {
	e1 := errors.New(1000, "error1", errors.Business)
	assert.Equal(t, 1000, e1.Code)
	assert.Equal(t, "error1", e1.Msg)
	assert.Equal(t, errors.Business, e1.Typ)

	e2 := errors.New(1000, "error1", errors.Framework)
	assert.Equal(t, 1000, e2.Code)
	assert.Equal(t, "error1", e2.Msg)
	assert.Equal(t, errors.Framework, e2.Typ)
}

func TestErrorType_String(t *testing.T) {
	et1 := errors.Framework
	assert.Equal(t, "framework error", et1.String())

	et2 := errors.Business
	assert.Equal(t, "business error", et2.String())
}

func TestError_Error(t *testing.T) {
	e1 := errors.New(1000, "error1", errors.Business)
	m1 := fmt.Sprintf("errCode: %d, errMsg: %s, errType: %s", e1.Code, e1.Msg, e1.Typ)
	assert.Equal(t, m1, e1.Error())
}
