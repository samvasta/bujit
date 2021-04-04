package actions_accounts

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

func TestDeleteAccountAction_Soft(t *testing.T) {
	session := session.InMemorySession(models.MigrateSchema)

	account := models.Account{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23)}}
	session.Db.Create(&account)

	action := DeleteAccountAction{Session: &session, Name: "Account 1", IsHardDelete: false}
	result, consequences := action.Execute()

	assert.True(t, result.IsSuccessful)

	assert.Len(t, consequences, 1)

	assert.Equal(t, actions.DELETE, consequences[0].ConsequenceType)
	assert.Equal(t, account, consequences[0].Object)

	// Make sure the account is still in the db and has state=closed
	var dbAccount models.Account
	session.Db.Joins("CurrentState").Find(&dbAccount, account.ID)

	assert.True(t, dbAccount.CurrentState.IsClosed)
}

func TestDeleteAccountAction_Hard(t *testing.T) {
	session := session.InMemorySession(models.MigrateSchema)

	originalState := models.AccountState{Balance: models.MakeMoney(0)}
	session.Db.Create(&originalState)

	account := models.Account{Name: "Account 1", Description: "description1", IsActive: true, CurrentState: models.AccountState{Balance: models.MakeMoney(1.23), PrevStateID: &originalState.ID}}
	session.Db.Create(&account)

	var numAccounts, numAccountStates int64
	session.Db.Model(&models.Account{}).Count(&numAccounts)
	session.Db.Model(&models.AccountState{}).Count(&numAccountStates)

	assert.Equal(t, int64(1), numAccounts)
	assert.Equal(t, int64(2), numAccountStates)

	action := DeleteAccountAction{Session: &session, Name: "Account 1", IsHardDelete: true}
	result, consequences := action.Execute()

	assert.True(t, result.IsSuccessful)

	assert.Len(t, consequences, 1)

	assert.Equal(t, actions.DELETE, consequences[0].ConsequenceType)
	assert.Equal(t, account, consequences[0].Object)

	// Make sure the account is still in the db and has state=closed

	session.Db.Model(&models.Account{}).Count(&numAccounts)
	session.Db.Model(&models.AccountState{}).Count(&numAccountStates)
	assert.Zero(t, numAccounts)
	assert.Zero(t, numAccountStates)
}
