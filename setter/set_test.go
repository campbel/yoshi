package setter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDuration(t *testing.T) {
	assert := assert.New(t)

	var d time.Duration
	assert.NoError(SetAny(&d, "1h"))
	assert.Equal(time.Hour, d)
}

func TestString(t *testing.T) {
	assert := assert.New(t)

	var s string
	assert.NoError(SetAny(&s, "foo"))
	assert.Equal("foo", s)
}
