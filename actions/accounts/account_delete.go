package actions_accounts

import (
	"fmt"

	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

type DeleteAccountAction struct {
	Name         string
	IsHardDelete bool // Soft delete closes the account. Hard delete permanently deletes all account data
	Session      *session.Session
}

func (action DeleteAccountAction) IsValid() bool {
	if action.Session == nil || action.Name == "" {
		return false
	}

	var accounts []models.Account
	result := action.Session.Db.Where("Name = ?", action.Name).Find(&accounts)

	return result.Error == nil && len(accounts) == 1
}

func (action DeleteAccountAction) Execute() (actions.ActionResult, []*actions.Consequence) {

	consequences := []*actions.Consequence{}

	var accounts []models.Account
	tx := action.Session.Db.Preload("CurrentState").Where("Name = ?", action.Name).Find(&accounts)

	if tx.Error != nil {
		return actions.ActionResult{Output: tx.Error.Error(), IsSuccessful: false}, []*actions.Consequence{}
	}

	if len(accounts) == 0 {
		// No account found matching name
		return actions.ActionResult{Output: fmt.Sprintf(`{"detail": "No account with name '%s'"}`, action.Name), IsSuccessful: false}, []*actions.Consequence{}
	}

	if action.IsHardDelete {
		for _, account := range accounts {
			consequences = append(consequences, &actions.Consequence{ConsequenceType: actions.DELETE, Object: account})
			//delete all account states
			deleteAccountState(account.CurrentState, action.Session)
			action.Session.Db.Delete(&account)
		}
	} else {
		for _, account := range accounts {
			consequences = append(consequences, &actions.Consequence{ConsequenceType: actions.DELETE, Object: account})

			currentState := account.CurrentState
			nextState := models.AccountState{
				Balance:     account.CurrentState.Balance,
				PrevState:   &currentState,
				PrevStateID: &currentState.ID,
				IsClosed:    true,
			}
			action.Session.Db.Create(&nextState)

			account.CurrentState = nextState
			account.CurrentStateID = &nextState.ID
			action.Session.Db.Save(&account)
		}
	}

	return actions.ActionResult{Output: "", IsSuccessful: true}, consequences

}

func deleteAccountState(state models.AccountState, session *session.Session) {
	session.Db.Preload("PrevState").First(&state, state.ID)
	prevStateId := state.PrevStateID
	session.Db.Delete(&state)
	if prevStateId != nil {
		prevState := *state.PrevState
		deleteAccountState(prevState, session)
	}
}
