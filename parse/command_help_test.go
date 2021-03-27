package parse

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"samvasta.com/bujit/session"
)

func TestHelpCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help"

	action := ParseExpression(input, session)

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
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help"

	action := ParseExpression(input, session)

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

func TestHelpPartialVerboseCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help --ver"

	action := ParseExpression(input, session)

	result, consequences := action.Execute()

	fmt.Println(result.Suggestions)

	assert.Len(t, consequences, 0)
	assert.Equal(t, "", result.Output)

	assert.Len(t, result.Suggestions, 1)
	assert.Contains(t, result.Suggestions, "--verbose")
}
