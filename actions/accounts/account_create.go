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

func (action CreateAccountAction) IsValid() bool {
	return action.Name != "" && action.Session != nil
}

func (action CreateAccountAction) Execute() (actions.ActionResult, []*actions.Consequence) {

	account := models.Account{Name: action.Name, Description: action.Description, IsActive: true, CurrentState: models.AccountState{Balance: action.StartingBalance}}
	tx := action.Session.Db.Create(&account)

	if tx.Error != nil {
		return actions.ActionResult{Output: tx.Error.Error(), IsSuccessful: false}, []*actions.Consequence{}
	}

	consequences := []*actions.Consequence{
		{ConsequenceType: actions.CREATE, Object: account},
	}

	if action.CategoryName != "" {
		// Upsert category
		var category models.Category
		tx := action.Session.Db.FirstOrCreate(&category, models.Category{Name: action.CategoryName})
		isCategoryNew := tx.RowsAffected > 0

		if tx.Error != nil {
			return actions.ActionResult{Output: tx.Error.Error(), IsSuccessful: false}, []*actions.Consequence{}
		}
		// Create association for account & category
		assocErr := action.Session.Db.Model(&category).Association("Accounts").Append(&account)

		if assocErr != nil {
			return actions.ActionResult{Output: tx.Error.Error(), IsSuccessful: false}, []*actions.Consequence{}
		}

		var categoryConsequenceType actions.ConsequenceType
		if isCategoryNew {
			categoryConsequenceType = actions.CREATE
		} else {
			categoryConsequenceType = actions.UPDATE
		}
		consequences = append(consequences,
			&actions.Consequence{ConsequenceType: categoryConsequenceType, Object: category})
	}

	return actions.ActionResult{Output: "", IsSuccessful: true}, consequences

}
