package actions_accounts

import (
	"fmt"
	"strings"

	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models"
	"samvasta.com/bujit/session"
)

type ListAccountAction struct {
	Name         string
	Description  string
	MinBalance   *models.Money
	MaxBalance   *models.Money
	CategoryName string
	AsTree       bool
	Session      *session.Session
}

type ListAccountOutput struct {
	Tree bool `json="tree"`
}

func (action ListAccountAction) IsValid() bool {
	return action.Session != nil
}

func (action ListAccountAction) Execute() (actions.ActionResult, []*actions.Consequence) {

	consequences := []*actions.Consequence{}

	conditions := []string{}
	var conditionValues []interface{}

	if action.Name != "" {
		conditions = append(conditions, "accounts.Name LIKE ?")
		conditionValues = append(conditionValues, "%"+action.Name+"%")
	}

	if action.Description != "" {
		conditions = append(conditions, "accounts.Description LIKE ?")
		conditionValues = append(conditionValues, "%"+action.Description+"%")
	}

	if action.CategoryName != "" {
		// conditions = append(conditions, "Category.Name <> ''")
		conditions = append(conditions, "Category.fully_qualified_name LIKE ?")
		conditionValues = append(conditionValues, "%"+action.CategoryName+"%")
	}

	if action.MinBalance != nil {
		value := (*action.MinBalance).Value()
		conditions = append(conditions, "CurrentState.Balance >= ?")
		conditionValues = append(conditionValues, fmt.Sprint(value))
	}

	if action.MaxBalance != nil {
		value := (*action.MaxBalance).Value()
		conditions = append(conditions, "CurrentState.Balance <= ?")
		conditionValues = append(conditionValues, fmt.Sprint(value))
	}

	var accounts []models.Account
	query := strings.Join(conditions, " AND ")
	action.Session.Db.Joins("CurrentState").Joins("Category").Where(query, conditionValues...).Find(&accounts)

	for _, a := range accounts {
		a.Session = action.Session
		consequences = append(consequences, &actions.Consequence{ConsequenceType: actions.READ, Object: a})
	}

	output := ListAccountOutput{Tree: action.AsTree}

	return actions.ActionResult{Output: output, IsSuccessful: true}, consequences

}
