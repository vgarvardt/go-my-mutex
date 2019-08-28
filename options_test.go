package mymutex

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWithTimeout(t *testing.T) {
	m, err := New(nil, WithTimeout(1234))
	require.NoError(t, err)
	assert.Equal(t, 1234, m.timeout)
}
