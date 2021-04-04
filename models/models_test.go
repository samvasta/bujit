package models

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"samvasta.com/bujit/session"
)

func TestCategoryMarshal(t *testing.T) {
	session := session.InMemorySession(MigrateSchema)
	session.CurrencyPrefix = ""
	session.CurrencySuffix = "USD"

	categoryId := uint(1)
	parentId := uint(2)

	state1 := AccountState{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(432.12),
		nil, nil,
		false,
		&session}

	state2 := AccountState{
		124,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(123.45),
		nil, nil,
		false,
		&session}

	account := Account{
		ID:             1,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:           "Account",
		Description:    "Account Description",
		IsActive:       true,
		CurrentStateID: &state1.ID,
		CurrentState:   state1,
		Session:        &session,
	}

	account2 := Account{
		ID:             2,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:           "Account2",
		Description:    "Account2 Description",
		IsActive:       true,
		CurrentStateID: &state2.ID,
		CurrentState:   state2,
		Session:        &session,
	}

	subCategory := Category{
		ID:                 3,
		CreatedAt:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		UpdatedAt:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:               "Test Category",
		FullyQualifiedName: "Test Category",
		Description:        "Description",
		SuperCategoryID:    &categoryId,
		SubCategories:      []Category{},
		Accounts:           []Account{account2}, // Make sure the currentBalance includes the balance of this account
		Session:            &session,
	}

	category := Category{
		ID:                 1,
		CreatedAt:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		UpdatedAt:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:               "Test Category",
		FullyQualifiedName: "Test Category",
		Description:        "Description",
		SuperCategoryID:    &parentId,
		SubCategories:      []Category{subCategory},
		Accounts:           []Account{account},
		Session:            &session,
	}

	jsonBytes, error := json.Marshal(category)
	jsonStr := string(jsonBytes)

	if error != nil {
		assert.Fail(t, "%v", error)
	}

	expected := `{
		"id": 1,
		"createdAt": "2020-01-01T00:00:00Z",
		"updatedAt": "2020-01-01T00:00:00Z",
		"name": "Test Category",
		"description": "Description",
		"accounts": [1],
		"currentBalance": "555.57 USD",
		"subCategories": [3]
	}`

	assert.JSONEq(t, expected, jsonStr)
}

func TestAccountStateMarshal(t *testing.T) {
	session := session.InMemorySession(MigrateSchema)
	session.CurrencyPrefix = ""
	session.CurrencySuffix = "USD"
	as := AccountState{
		123,
		time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		MakeMoney(432.12),
		nil, nil,
		false,
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
		false,
		&session}
	account := Account{
		ID:             123,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:           "Account Name",
		Description:    "description",
		IsActive:       false,
		CurrentStateID: &state.ID,
		CurrentState:   state,
		Session:        &session}

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
		false,
		&session}
	fromAccount := Account{
		ID:             123,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:           "Source",
		Description:    "description",
		IsActive:       false,
		CurrentStateID: &fromAccountState.ID,
		CurrentState:   fromAccountState,
		Session:        &session}

	toAccountState := AccountState{
		2,
		time.Now().Unix(),
		MakeMoney(12.34),
		nil, nil,
		false,
		&session}
	toAccount := Account{
		ID:             123,
		CreatedAt:      time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Name:           "Sink",
		Description:    "description",
		IsActive:       false,
		CurrentStateID: &toAccountState.ID,
		CurrentState:   toAccountState,
		Session:        &session}

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
