package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestAccountListCommand(t *testing.T) {
	session := session.InMemorySession(models.MigrateSchema)

	category1 := models.MakeCategory("cat1", "", nil)
	category2 := models.MakeCategory("other", "", nil)
	accounts := []*models.Account{
		{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23)}},
		{Name: "Account 2", Description: "description2", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(4.56)}, CategoryID: &category1.ID, Category: category1},
		{Name: "Account 3", Description: "description3", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(7.89)}, CategoryID: &category2.ID, Category: category2},
	}

	category1.Accounts = append(category1.Accounts, *accounts[1])
	category2.Accounts = append(category2.Accounts, *accounts[2])

	for _, a := range accounts {
		session.Db.Create(a)
	}

	testCase := func(input string, isValid bool, expectedSuggestions []string, additionalCheck func(t *testing.T, action actions.Actioner)) func(t *testing.T) {
		return func(t *testing.T) {

			action, suggestion := ParseExpression(input, session)

			assert.Equal(t, isValid, suggestion.isValidAsIs)

			for _, expectedSuggestion := range expectedSuggestions {
				assert.Contains(t, suggestion.nextArgs, expectedSuggestion)
			}
			assert.Len(t, suggestion.nextArgs, len(expectedSuggestions))

			additionalCheck(t, action)
		}
	}

	t.Run("list account",
		testCase("list account",
			true,
			[]string{"--name", "--description", "--category", "--max-balance", "--min-balance", "--help"},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)
			}))

	t.Run("fully specified",
		testCase("list account -n name -d description --category='category' -m $123.45 -x 432.11",
			true,
			[]string{},
			func(test *testing.T, action actions.Actioner) {
				assert.NotNil(t, action)

				listAccountAction := action.(actions_accounts.ListAccountAction)

				assert.Equal(t, "name", listAccountAction.Name)
				assert.Equal(t, "description", listAccountAction.Description)
				assert.Equal(t, "category", listAccountAction.CategoryName)
				assert.Equal(t, models.MakeMoney(123.45), *listAccountAction.MinBalance)
				assert.Equal(t, models.MakeMoney(432.11), *listAccountAction.MaxBalance)
			}))
}
