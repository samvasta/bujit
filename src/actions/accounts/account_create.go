package actions_accounts

import (
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

type CreateAccountAction struct {
	Name            string
	Description     string
	StartingBalance models.Money
	CategoryName    string
	Session         *session.Session
}

func (action *CreateAccountAction) execute() (actions.ActionResult, []*actions.Consequence) {

	// Upsert category
	var category models.Category
	tx := action.Session.Db.FirstOrCreate(&category, models.Category{Name: action.CategoryName})
	isCategoryNew := tx.RowsAffected > 0

	if tx.Error != nil {
		return actions.ActionResult{Output: tx.Error.Error()}, []*actions.Consequence{}
	}

	account := models.Account{Name: action.Name, Description: action.Description, IsActive: true, CurrentState: models.AccountState{Balance: action.StartingBalance}}
	tx = action.Session.Db.Create(&account)

	if tx.Error != nil {
		return actions.ActionResult{Output: tx.Error.Error()}, []*actions.Consequence{}
	}

	// Create association for account & category
	assocErr := action.Session.Db.Model(&category).Association("Accounts").Append(&account)

	if assocErr != nil {
		return actions.ActionResult{Output: tx.Error.Error()}, []*actions.Consequence{}
	}

	var categoryConsequenceType actions.ConsequenceType
	if isCategoryNew {
		categoryConsequenceType = actions.CREATE
	} else {
		categoryConsequenceType = actions.UPDATE
	}
	return actions.ActionResult{Output: ""},
		[]*actions.Consequence{
			{ConsequenceType: actions.CREATE, Object: account},
			{ConsequenceType: categoryConsequenceType, Object: category},
		}
}
