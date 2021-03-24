package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/session"
)

func TestAccountStateMarshal(t *testing.T) {
	session := session.InMemorySession(MigrateSchema)
	session.CurrencyPrefix = ""
	session.CurrencySuffix = "USD"
	as := AccountState{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(432.12),
		nil, nil,
		&session}

	jsonBytes, error := json.Marshal(as)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"balance":"432.12 USD",
		"timestamp":"2020-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expected, jsonStr)
}

func TestAccountMarshal(t *testing.T) {
	session := session.InMemorySession(MigrateSchema)
	session.CurrencyPrefix = ""
	session.CurrencySuffix = "USD"
	state := AccountState{
		123,
		time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(123.45),
		nil, nil,
		&session}
	account := Account{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		"Account Name",
		"description",
		false,
		&state.ID,
		state,
		&session}

	jsonBytes, error := json.Marshal(account)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"createdAt":"2020-01-01T00:00:00Z",
		"currentBalance":"123.45 USD",
		"description":"description",
		"name":"Account Name",
		"status":"closed",
		"updatedAt":"2021-01-01T00:00:00Z"
		}`

	assert.JSONEq(t, expected, jsonStr)

	account.IsActive = true
	jsonBytes, error = json.Marshal(account)
	jsonStr = string(jsonBytes)

	assert.JSONEq(t, strings.Replace(expected, "closed", "open", 1), jsonStr)
}

func TestTransactionMarshal(t *testing.T) {
	session := session.InMemorySession(MigrateSchema)
	session.CurrencyPrefix = ""
	session.CurrencySuffix = "USD"

	fromAccountState := AccountState{
		1,
		time.Now().Unix(),
		MakeMoney(12.34),
		nil, nil,
		&session}
	fromAccount := Account{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		"Source",
		"description",
		false,
		&fromAccountState.ID,
		fromAccountState,
		&session}

	toAccountState := AccountState{
		2,
		time.Now().Unix(),
		MakeMoney(12.34),
		nil, nil,
		&session}
	toAccount := Account{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		"Sink",
		"description",
		false,
		&toAccountState.ID,
		toAccountState,
		&session}

	transaction := Transaction{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(123.45),
		&fromAccount.ID,
		&fromAccount,
		&toAccount.ID,
		&toAccount,
		"Memo",
		&session}

	jsonBytes, error := json.Marshal(transaction)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"memo":"Memo",
		"amount":"123.45 USD",
		"fromAccount":"Source",
		"toAccount":"Sink",
		"timestamp":"2020-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expected, jsonStr)

}
