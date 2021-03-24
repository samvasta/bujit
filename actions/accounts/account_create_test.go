package actions_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestCreateAccount(t *testing.T) {
	session := session.InMemorySession(models.MigrateSchema)

	action := CreateAccountAction{Name: "Account Name", Description: "Account Description", StartingBalance: models.MakeMoney(123.45), CategoryName: "New Category", Session: &session}

	result, consequences := action.Execute()

	// Test that database has correct data

	var account models.Account
	var accountState models.AccountState
	var category models.Category
	accountResult := session.Db.Preload("CurrentState").Find(&account, models.Account{Name: "Account Name"})
	accountStateResult := session.Db.Find(&accountState, models.AccountState{ID: *account.CurrentStateID})
	categoryResult := session.Db.Preload("Accounts.CurrentState").Find(&category, models.Category{Name: "New Category"})

	assert.Nil(t, accountResult.Error)
	assert.Nil(t, accountStateResult.Error)
	assert.Nil(t, categoryResult.Error)

	// Account State
	assert.Equal(t, accountState, account.CurrentState)
	assert.Nil(t, accountState.PrevStateID)
	assert.Equal(t, int64(12345), accountState.Balance.Value())

	// Account
	assert.Equal(t, "Account Name", account.Name)
	assert.Equal(t, models.MakeMoney(123.45), account.Balance())
	assert.Equal(t, "Account Description", account.Description)
	assert.True(t, account.IsActive)

	// Category
	assert.Len(t, category.Accounts, 1, "new category should only have one account")
	assert.Equal(t, account, category.Accounts[0])
	assert.Empty(t, category.SubCategories)
	assert.Equal(t, "New Category", category.Name)
	assert.Nil(t, category.SuperCategoryID, "new category should not belong to another category")

	// Test that return values are correct

	assert.Len(t, consequences, 2)
	assert.Empty(t, result.Output)
	accountCreateConsequence := consequences[0]
	categoryCreateConsequence := consequences[1]

	assert.Equal(t, actions.CREATE, accountCreateConsequence.ConsequenceType)
	assert.Equal(t, account, accountCreateConsequence.Object)

	assert.Equal(t, actions.CREATE, categoryCreateConsequence.ConsequenceType)
	// assert.Equal(t, category, categoryCreateConsequence.Object)
}
