package parse

import (
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
	_, consequences := action.Execute()

	assert.True(t, suggestion.IsValidAsIs)
	assert.Empty(t, suggestion.NextArgs)

	assert.Len(t, consequences, 0)
}

func TestHelpVerboseCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help"

	action, suggestion := ParseExpression(input, &session)

	assert.NotNil(t, action)
	_, consequences := action.Execute()

	assert.True(t, suggestion.IsValidAsIs)
	assert.Empty(t, suggestion.NextArgs)
	assert.Len(t, consequences, 0)
}

func TestHelpPartialVerboseCommand(t *testing.T) {
	session := session.InMemorySession(func(d *gorm.DB) {})
	input := "help --ver"

	action, suggestion := ParseExpression(input, &session)

	assert.Nil(t, action)

	assert.Len(t, suggestion.NextArgs, 1)
	assert.Contains(t, suggestion.NextArgs, "--verbose")
}
