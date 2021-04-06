package actions

import (
	"encoding/json"
)

type ActionType int

const (
	// Meta Actions
	EXIT ActionType = iota
	VERSION
	HELP
	SUGGEST

	// Account Actions
	NEW_ACCOUNT
	DELETE_ACCOUNT
	MODIFY_ACCOUNT
	LIST_ACCOUNT

	// Account State Actions
	LIST_ACCOUNT_STATE

	// Transaction Actions
	DETAIL_ACCOUNT
	NEW_TRANSACTION
	LIST_TRANSACTION
)

type ConsequenceType int

const (
	CREATE ConsequenceType = iota
	READ
	UPDATE
	DELETE
)

func (this ConsequenceType) String() string {
	switch this {
	case CREATE:
		return "Create"
	case READ:
		return "Read"
	case UPDATE:
		return "Update"
	case DELETE:
		return "Delete"
	default:
		return ""
	}
}

type Consequence struct {
	ConsequenceType ConsequenceType
	Object          json.Marshaler
}

type ActionResult struct {
	Output       interface{}
	IsSuccessful bool
}

type Actioner interface {
	Execute() (ActionResult, []*Consequence)
}
