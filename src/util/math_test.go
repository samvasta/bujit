package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsI(t *testing.T) {
	assert.Equal(t, 20, AbsI(20))
	assert.Equal(t, 20, AbsI(-20))
}

func TestAbsI64(t *testing.T) {
	assert.Equal(t, int64(20), AbsI64(int64(20)))
	assert.Equal(t, int64(20), AbsI64(int64(-20)))
}
