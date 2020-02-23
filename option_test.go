package gorpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithConfig(t *testing.T) {
	opts := options{}
	o := WithConfig("../testcase/service.ini")
	o(&opts)

	assert.Equal(t, "../testcase/service.ini", opts.conf)
}
