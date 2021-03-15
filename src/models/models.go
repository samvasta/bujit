package models

import (
	"time"

	humanize "github.com/dustin/go-humanize"
)

type Detailser interface {
	Details() map[string]interface{}
}

type AccountState struct {
	id         int64
	timestamp  time.Time
	balance    *Money
	prev_state *AccountState
}

func (accountState *AccountState) Details() (details map[string]interface{}) {
	details["ID"] = accountState.id
	details["Timestamp"] = humanize.Time(accountState.timestamp)
	details["Balance"] = accountState.balance.String()

	return details
}

type Account struct {
	id            int64
	name          string
	isActive      bool
	description   string
	created_on    time.Time
	current_state *AccountState
}

func (account *Account) Balance() *Money {
	return account.current_state.balance
}

func (account *Account) Details() (details map[string]interface{}) {
	details["ID"] = account.id
	details["Name"] = account.name

	if account.isActive {
		details["Status"] = "open"
	} else {
		details["Status"] = "closed"
	}

	details["Description"] = account.description
	details["Created on"] = humanize.Time(account.created_on)
	details["Updated on"] = humanize.Time(account.current_state.timestamp)
	details["Current Balance"] = account.Balance().String()

	return details
}

type Transaction struct {
	id          int64
	timestamp   time.Time
	change      Money
	source      *Account
	destination *Account
	memo        string
}

func (tran *Transaction) SourceExists() bool {
	return tran.source != nil
}
func (tran *Transaction) DestinationExists() bool {
	return tran.destination != nil
}

func (tran *Transaction) Details() (details map[string]interface{}) {
	details["ID"] = tran.id
	details["Timestamp"] = humanize.Time(tran.timestamp)
	details["Amount"] = tran.change.String()

	if tran.SourceExists() {
		details["Source Account"] = tran.source.name
	}
	if tran.DestinationExists() {
		details["Destination Account"] = tran.destination.name
	}

	details["Memo"] = tran.memo

	return details
}
