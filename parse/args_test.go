package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestParseOptionalArg(t *testing.T) {
	argTok := MakeOptionalArgToken(1, "o", "option")

	testCase := func(input, output string, isValid bool, suggestions []string) func(t *testing.T) {
		return func(t *testing.T) {

			session := session.InMemorySession(models.MigrateSchema)

			ctx := ParseContext{tokens: Tokenize(input), currentTokenIndex: 0, isValid: true, session: &session}
			value, suggestion := parseOptionalArg(&ctx, argTok, ItemNamePattern, "patternName")

			nextToken, hasNext := ctx.nextToken()

			assert.Equal(t, "", nextToken)
			assert.False(t, hasNext)

			assert.Equal(t, output, value)
			assert.Equal(t, isValid, suggestion.IsValidAsIs)

			for _, expected := range suggestions {
				assert.Contains(t, suggestion.NextArgs, expected)
			}
		}
	}

	t.Run("invalid - missing arg", testCase("-o", "", false, []string{"<patternName>"}))
	t.Run("invalid - no input", testCase("", "", false, []string{"--option"}))

	t.Run("valid - short name", testCase("-o value", "value", true, []string{}))
	t.Run("valid - short name with equals", testCase("-o=value", "value", true, []string{}))
	t.Run("valid - long name", testCase("--option value", "value", true, []string{}))
	t.Run("valid - long name with equals", testCase("--option=value", "value", true, []string{}))

	t.Run("valid - short name with quotes", testCase("-o 'value'", "value", true, []string{}))
	t.Run("valid - short name with equals and quotes", testCase("-o=\"value\"", "value", true, []string{}))
	t.Run("valid - long name with quotes", testCase("--option \"value\"", "value", true, []string{}))
	t.Run("valid - long name with equals and quotes", testCase("--option='value'", "value", true, []string{}))
}
