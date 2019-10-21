package pool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithOptions(t *testing.T) {
	opts := &Options{}
	WithMinIdle(1)(opts)
	WithMaxIdle(2)(opts)
	WithMaxActive(10)(opts)
	WithIdleTimeout(time.Second)(opts)
	WithDialTimeout(time.Second)(opts)
	WithMaxConnLifetime(time.Second * 60)(opts)
	WithWait(true)(opts)

	assert.Equal(t, opts.MinIdle, 1)
	assert.Equal(t, opts.MaxIdle, 2)
	assert.Equal(t, opts.MaxActive, 10)
	assert.Equal(t, opts.IdleTimeout, time.Second)
	assert.Equal(t, opts.DialTimeout, time.Second)
	assert.Equal(t, opts.MaxConnLifetime, 60*time.Second)
	assert.Equal(t, opts.Wait, true)
}
