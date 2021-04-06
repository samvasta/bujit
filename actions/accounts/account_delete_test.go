package actions_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestDeleteAccountAction_Soft(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	account := models.Account{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23)}, Session: &s}
	s.Db.Create(&account)

	action := DeleteAccountAction{Session: &s, Name: "Account 1", IsHardDelete: false}
	result, consequences := action.Execute()

	assert.True(t, result.IsSuccessful)

	assert.Len(t, consequences, 1)

	assert.Equal(t, actions.DELETE, consequences[0].ConsequenceType)
	assert.Equal(t, account, consequences[0].Object)

	// Make sure the account is still in the db and has state=closed
	var dbAccount models.Account
	s.Db.Joins("CurrentState").Find(&dbAccount, account.ID)

	assert.True(t, dbAccount.CurrentState.IsClosed)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}

func TestDeleteAccountAction_Hard(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	originalState := models.AccountState{Balance: models.MakeMoney(0), Session: &s}
	s.Db.Create(&originalState)

	account := models.Account{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23), PrevStateID: &originalState.ID}, Session: &s}
	s.Db.Create(&account)

	var numAccounts, numAccountStates int64
	s.Db.Model(&models.Account{}).Count(&numAccounts)
	s.Db.Model(&models.AccountState{}).Count(&numAccountStates)

	assert.Equal(t, int64(1), numAccounts)
	assert.Equal(t, int64(2), numAccountStates)

	action := DeleteAccountAction{Session: &s, Name: "Account 1", IsHardDelete: true}
	result, consequences := action.Execute()

	assert.True(t, result.IsSuccessful)

	assert.Len(t, consequences, 1)

	assert.Equal(t, actions.DELETE, consequences[0].ConsequenceType)
	assert.Equal(t, account, consequences[0].Object)

	// Make sure the account is still in the db and has state=closed

	s.Db.Model(&models.Account{}).Count(&numAccounts)
	s.Db.Model(&models.AccountState{}).Count(&numAccountStates)
	assert.Zero(t, numAccounts)
	assert.Zero(t, numAccountStates)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}

func TestDeleteAccountAction_DoesNotExist(t *testing.T) {
	s := session.InMemorySession(models.MigrateSchema)

	account := models.Account{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23)}}
	s.Db.Create(&account)

	action := DeleteAccountAction{Session: &s, Name: "Does not exist", IsHardDelete: false}
	result, consequences := action.Execute()

	assert.False(t, result.IsSuccessful)
	assert.JSONEq(t, `{"detail": "No account with name 'Does not exist'"}`, result.Output)
	assert.Len(t, consequences, 0)

	for _, c := range consequences {
		assert.Equal(t, s, *c.Object.(session.Sessioner).GetSession())
	}
}
