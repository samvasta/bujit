package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestAccountCreateCommand_new_account(t *testing.T) {
	testCase := func(input, output string, isValid bool, expectedSuggestions []string, additionalCheck func(t *testing.T, action actions.Actioner)) func(t *testing.T) {
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

	t.Run("new account",
		testCase("new account",
			"",
			false,
			[]string{"<name>", "--help"},
			func(test *testing.T, action actions.Actioner) {
				assert.Nil(t, action)
			}))

	t.Run("new account name",
		testCase("new account name",
			"",
			true,
			[]string{"--description", "--category", "--balance"},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				createAccountAction := action.(actions_accounts.CreateAccountAction)

				assert.Equal(t, "name", createAccountAction.Name)
				assert.Equal(t, "", createAccountAction.CategoryName)
				assert.Equal(t, "", createAccountAction.Description)
				assert.Equal(t, models.MakeMoney(0), createAccountAction.StartingBalance)
			}))

	t.Run("new account name description",
		testCase("new account name -d='description'",
			"",
			true,
			[]string{"--category", "--balance"},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				createAccountAction := action.(actions_accounts.CreateAccountAction)

				assert.Equal(t, "name", createAccountAction.Name)
				assert.Equal(t, "", createAccountAction.CategoryName)
				assert.Equal(t, "description", createAccountAction.Description)
				assert.Equal(t, models.MakeMoney(0), createAccountAction.StartingBalance)
			}))

	t.Run("new account name description without arg value",
		testCase("new account name -d=",
			"",
			false,
			[]string{"<description>"},
			func(test *testing.T, action actions.Actioner) {
				assert.Nil(t, action)
			}))

	t.Run("new account name category description",
		testCase("new account name --category \"Test Category\" -d='description'",
			"",
			true,
			[]string{"--balance"},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				createAccountAction := action.(actions_accounts.CreateAccountAction)

				assert.Equal(t, "name", createAccountAction.Name)
				assert.Equal(t, "Test Category", createAccountAction.CategoryName)
				assert.Equal(t, "description", createAccountAction.Description)
				assert.Equal(t, models.MakeMoney(0), createAccountAction.StartingBalance)
			}))

	t.Run("new account name category balance description",
		testCase("new account name --category \"Test Category\" -b $1.23 -d='description'",
			"",
			true,
			[]string{},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				createAccountAction := action.(actions_accounts.CreateAccountAction)

				assert.Equal(t, "name", createAccountAction.Name)
				assert.Equal(t, "Test Category", createAccountAction.CategoryName)
				assert.Equal(t, "description", createAccountAction.Description)
				assert.Equal(t, models.MakeMoney(1.23), createAccountAction.StartingBalance)
			}))
}
