package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func buildTestCase(input, output string, isValid bool, expectedSuggestions []string, additionalCheck func(t *testing.T, action actions.Actioner)) func(t *testing.T) {
	return func(t *testing.T) {
		session := session.InMemorySession(models.MigrateSchema)

		action, suggestion := ParseExpression(input, session)

		assert.Equal(t, isValid, suggestion.isValidAsIs)

		for _, expectedSuggestion := range expectedSuggestions {
			assert.Contains(t, suggestion.nextArgs, expectedSuggestion)
		}
		assert.Len(t, suggestion.nextArgs, len(expectedSuggestions))

		additionalCheck(t, action)
	}
}

func TestAccountCreateCommand_new_account(t *testing.T) {
	t.Run("new account",
		buildTestCase("new account",
			"",
			false,
			[]string{"<name>", "--help"},
			func(test *testing.T, action actions.Actioner) {
				assert.Nil(t, action)
			}))

	t.Run("new account name",
		buildTestCase("new account name",
			"",
			true,
			[]string{"--description=<STR>", "--category=<STR>", "--balance=<AMOUNT>"},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				createAccountAction := action.(actions_accounts.CreateAccountAction)

				assert.Equal(t, "name", createAccountAction.Name)
				assert.Equal(t, "", createAccountAction.CategoryName)
				assert.Equal(t, "", createAccountAction.Description)
				assert.Equal(t, models.MakeMoney(0), createAccountAction.StartingBalance)
			}))
}
