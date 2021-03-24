package parse

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelpCommand(t *testing.T) {
	input := "help"

	action := ParseExpression(input)

	result, consequences := action.Execute()

	assert.Len(t, consequences, 0)

	var unmarshaled []map[string]interface{}
	err := json.Unmarshal([]byte(result.Output), &unmarshaled)

	assert.Nil(t, err)

	for _, data := range unmarshaled {
		assert.Contains(t, data, "data")

		help := data["data"]
		assert.NotNil(t, help)
	}
}

func TestHelpVerboseCommand(t *testing.T) {
	input := "help"

	action := ParseExpression(input)

	result, consequences := action.Execute()

	assert.Len(t, consequences, 0)

	var unmarshaled []map[string]interface{}
	err := json.Unmarshal([]byte(result.Output), &unmarshaled)

	assert.Nil(t, err)

	for _, data := range unmarshaled {
		assert.Contains(t, data, "data")

		help := data["data"]
		assert.NotNil(t, help)
	}
}
