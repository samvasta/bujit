package actions_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestListAccountAction(t *testing.T) {
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

	testCase := func(action ListAccountAction, expected []models.Account) func(t *testing.T) {
		return func(t *testing.T) {
			result, consequences := action.Execute()

			assert.True(t, result.IsSuccessful)
			if action.AsTree {
				assert.JSONEq(t, `{"tree": true}`, result.Output)
			} else {
				assert.JSONEq(t, `{"tree": false}`, result.Output)
			}

			assert.Len(t, consequences, len(expected))

			var actual []models.Account
			for _, c := range consequences {
				actual = append(actual, c.Object.(models.Account))

				assert.Equal(t, actions.READ, c.ConsequenceType)
			}
			assert.ElementsMatch(t, expected, actual)

		}
	}

	fourDollars := models.MakeMoney(4.0)
	fiveDollars := models.MakeMoney(5.0)

	t.Run("no args", testCase(ListAccountAction{Session: &session}, []models.Account{*accounts[0], *accounts[1], *accounts[2]}))
	t.Run("Name like '1'", testCase(ListAccountAction{Name: "1", Session: &session}, []models.Account{*accounts[0]}))
	t.Run("Description like '2'", testCase(ListAccountAction{Description: "2", Session: &session}, []models.Account{*accounts[1]}))
	t.Run("Category like 'cat'", testCase(ListAccountAction{CategoryName: "cat", Session: &session}, []models.Account{*accounts[1]}))
	t.Run("balance < 5", testCase(ListAccountAction{MaxBalance: &fiveDollars, Session: &session}, []models.Account{*accounts[0], *accounts[1]}))
	t.Run("4 < balance < 5", testCase(ListAccountAction{MinBalance: &fourDollars, MaxBalance: &fiveDollars, Session: &session}, []models.Account{*accounts[1]}))

	t.Run("as tree=true", testCase(ListAccountAction{Session: &session, AsTree: true}, []models.Account{*accounts[0], *accounts[1], *accounts[2]}))

}
