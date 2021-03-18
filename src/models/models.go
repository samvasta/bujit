package models

import (
	"encoding/json"
	"time"
)

type AccountState struct {
	Id        int64         `json:"id"`
	Timestamp time.Time     `json:"timestamp"`
	Balance   Money         `json:"balance"`
	prevState *AccountState `json:"-"`
}

type Account struct {
	Id           int64
	Name         string
	IsActive     bool
	Description  string
	CreatedOn    time.Time
	CurrentState *AccountState
}

func (account *Account) Balance() Money {
	return account.CurrentState.Balance
}

func (account Account) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = account.Id
	details["name"] = account.Name

	if account.IsActive {
		details["status"] = "open"
	} else {
		details["status"] = "closed"
	}

	details["description"] = account.Description
	details["createdOn"] = account.CreatedOn
	details["updatedOn"] = account.CurrentState.Timestamp
	details["currentBalance"] = account.Balance().String()

	return json.Marshal(details)
}

type Transaction struct {
	Id          int64
	Timestamp   time.Time
	Change      Money
	Source      *Account
	Destination *Account
	Memo        string
}

func (tran *Transaction) SourceExists() bool {
	return tran.Source != nil
}
func (tran *Transaction) DestinationExists() bool {
	return tran.Destination != nil
}

func (tran Transaction) MarshalJSON() ([]byte, error) {
	details := make(map[string]interface{})
	details["id"] = tran.Id
	details["timestamp"] = tran.Timestamp
	details["amount"] = tran.Change.String()

	if tran.SourceExists() {
		details["fromAccount"] = tran.Source.Name
	}
	if tran.DestinationExists() {
		details["toAccount"] = tran.Destination.Name
	}

	details["memo"] = tran.Memo

	return json.Marshal(details)
}
