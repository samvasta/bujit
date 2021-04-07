package outputview

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrappedString(t *testing.T) {
	input := "hello, world! Goodbye friendly world. thisistoolongtofitononeline."

	output := WrappedString(input, 0, 0, 20, 1000)

	expected :=
		`hello, world!
Goodbye friendly
world.
thisistoolongtofiton
oneline.
`
	assert.Equal(t, expected, output)
}

func TestWrappedString_Indented(t *testing.T) {
	input := "hello, world! Goodbye friendly world. thisistoolongtofitononeline."

	output := WrappedString(input, 0, 5, 25, 1000)

	expected :=
		`     hello, world!
     Goodbye friendly
     world.
     thisistoolongtofiton
     oneline.
`
	assert.Equal(t, expected, output)
}

func TestWrappedString_StartIndent(t *testing.T) {
	input := "hello, world! Goodbye friendly world. thisistoolongtofitononeline."

	output := WrappedString(input, 0, 5, 25, 1000)

	expected :=
		`hello, world! Goodbye
     friendly world.
     thisistoolongtofiton
     oneline.
`
	assert.Equal(t, expected, output)
}

func TestWrappedString_StartIndent2(t *testing.T) {
	input := "hello, world! Goodbye friendly world. thisistoolongtofitononeline."

	output := WrappedString(input, 15, 5, 25, 1000)

	expected :=
		`hello,
     world! Goodbye
     friendly world.
     thisistoolongtofiton
     oneline.
`
	assert.Equal(t, expected, output)
}
