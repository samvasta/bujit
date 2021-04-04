package parse

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"samvasta.com/bujit/session"
)

func TestHelpCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help"

	action, suggestion := ParseExpression(input, &session)

	assert.NotNil(t, action)
	result, consequences := action.Execute()

	assert.True(t, suggestion.IsValidAsIs)
	assert.Empty(t, suggestion.NextArgs)

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
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help"

	action, suggestion := ParseExpression(input, &session)

	assert.NotNil(t, action)
	result, consequences := action.Execute()

	assert.True(t, suggestion.IsValidAsIs)
	assert.Empty(t, suggestion.NextArgs)
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

func TestHelpPartialVerboseCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help --ver"

	action, suggestion := ParseExpression(input, &session)

	assert.Nil(t, action)

	assert.Len(t, suggestion.NextArgs, 1)
	assert.Contains(t, suggestion.NextArgs, "--verbose")
}
