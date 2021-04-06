package actions_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestCreateAccount(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	action := CreateAccountAction{Name: "Account Name", Description: "Account Description", StartingBalance: models.MakeMoney(123.45), CategoryName: "New Category", Session: &s}

	result, consequences := action.Execute()

	// Test that database has correct data

	var account models.Account
	var accountState models.AccountState
	var category models.Category
	accountResult := s.Db.Preload("CurrentState").Find(&account, models.Account{Name: "Account Name"})
	accountStateResult := s.Db.Find(&accountState, models.AccountState{ID: *account.CurrentStateID})
	categoryResult := s.Db.Preload("Accounts.CurrentState").Find(&category, models.Category{Name: "New Category"})

	// Set the expected session reference.
	account.Session = nil
	accountState.Session = nil // Should not be preloaded
	category.Session = nil     // Should not be preloaded

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
	assert.True(t, result.IsSuccessful)
	assert.Len(t, consequences, 2)
	assert.Empty(t, result.Output)
	accountCreateConsequence := consequences[0]
	categoryCreateConsequence := consequences[1]

	assert.Equal(t, actions.CREATE, accountCreateConsequence.ConsequenceType)
	account.Session = &s // fix expected session ptr
	assert.Equal(t, account, accountCreateConsequence.Object)

	assert.Equal(t, actions.CREATE, categoryCreateConsequence.ConsequenceType)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}

func TestCreateAccountEmptyCategoryNameDoesNotCreateCategory(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	action := CreateAccountAction{Name: "Account Name", Description: "Account Description", StartingBalance: models.MakeMoney(123.45), CategoryName: "", Session: &s}

	_, consequences := action.Execute()

	var numCategories int64
	s.Db.Model(&models.Category{}).Count(&numCategories)
	assert.Zero(t, numCategories)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}

func TestCreateAccountExistingCategory(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	var originalCategory models.Category
	s.Db.FirstOrCreate(&originalCategory, models.MakeCategory("Existing Category", "", nil))

	action := CreateAccountAction{Name: "Account Name", Description: "Account Description", StartingBalance: models.MakeMoney(123.45), CategoryName: "Existing Category", Session: &s}

	result, consequences := action.Execute()

	// Test that database has correct data

	var account models.Account
	var category models.Category
	accountResult := s.Db.Preload("CurrentState").Find(&account, models.Account{Name: "Account Name"})
	categoryResult := s.Db.Preload("Accounts.CurrentState").Find(&category, models.Category{Name: "Existing Category"})

	assert.Nil(t, accountResult.Error)
	assert.Nil(t, categoryResult.Error)

	// Category
	assert.Len(t, category.Accounts, 1, "new category should only have one account")
	assert.Equal(t, account, category.Accounts[0])
	assert.Empty(t, category.SubCategories)
	assert.Equal(t, "Existing Category", category.Name)
	assert.Nil(t, category.SuperCategoryID, "new category should not belong to another category")

	// Test that return values are correct
	assert.True(t, result.IsSuccessful)
	assert.Len(t, consequences, 2)
	assert.Empty(t, result.Output)
	accountCreateConsequence := consequences[0]
	categoryCreateConsequence := consequences[1]

	assert.Equal(t, actions.CREATE, accountCreateConsequence.ConsequenceType)
	account.Session = &s // fix the expected session ptr
	assert.Equal(t, account, accountCreateConsequence.Object)

	assert.Equal(t, actions.UPDATE, categoryCreateConsequence.ConsequenceType)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}

func TestCreateAccountDuplicateName(t *testing.T) {
	session := session.InMemorySession(models.MigrateSchema)

	existingAccount := models.Account{Name: "Account Name", Description: "description", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23)}}
	session.Db.Create(&existingAccount)

	action := CreateAccountAction{Name: "Account Name", Description: "Account Description", StartingBalance: models.MakeMoney(123.45), CategoryName: "Existing Category", Session: &session}

	result, _ := action.Execute()

	assert.False(t, result.IsSuccessful)
}
