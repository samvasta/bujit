package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCents(t *testing.T) {
	var m Money = 123

	cents := m.Cents()

	assert.Equal(t, int64(23), cents)
}

func TestDollars(t *testing.T) {
	var m Money = 1234

	dollars := m.Dollars()

	assert.Equal(t, int64(12), dollars)
}

func TestIsNegative(t *testing.T) {
	var m Money = -123

	assert.True(t, m.IsNegative())

	m = 123

	assert.False(t, m.IsNegative())
}

func TestValue(t *testing.T) {
	var m Money = -123

	assert.Equal(t, int64(-123), m.Value())
}
