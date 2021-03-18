package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAccountStateMarshal(t *testing.T) {
	as := AccountState{123, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), MakeMoney(432.12), nil}

	jsonBytes, error := json.Marshal(as)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"balance":"$432.12",
		"timestamp":"2020-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expected, jsonStr)
}

func TestAccountMarshal(t *testing.T) {
	state := AccountState{123, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), MakeMoney(123.45), nil}
	account := Account{123, "Account Name", false, "description", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), &state}

	jsonBytes, error := json.Marshal(account)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"createdOn":"2020-01-01T00:00:00Z",
		"currentBalance":"$123.45",
		"description":"description",
		"name":"Account Name",
		"status":"closed",
		"updatedOn":"2021-01-01T00:00:00Z"
		}`

	assert.JSONEq(t, expected, jsonStr)

	account.IsActive = true
	jsonBytes, error = json.Marshal(account)
	jsonStr = string(jsonBytes)

	assert.JSONEq(t, strings.Replace(expected, "closed", "open", 1), jsonStr)
}

func TestTransactionMarshal(t *testing.T) {
	fromAccount := Account{123, "Source", false, "description", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil}
	toAccount := Account{123, "Sink", false, "description", time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), nil}

	transaction := Transaction{123, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), MakeMoney(123.45), &fromAccount, &toAccount, "Memo"}

	jsonBytes, error := json.Marshal(transaction)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id":123,
		"memo":"Memo",
		"amount":"$123.45",
		"fromAccount":"Source",
		"toAccount":"Sink",
		"timestamp":"2020-01-01T00:00:00Z"
	}`

	assert.JSONEq(t, expected, jsonStr)

}
