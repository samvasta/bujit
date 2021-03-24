package actions

import (
	"strings"
)

type AutoSuggestAction struct {
	isValidAsIs bool
	nextArgs    []string
}

func (autoSuggestAction AutoSuggestAction) Execute() (result ActionResult, consequences []*Consequence) {

	var sb strings.Builder

	for i, arg := range autoSuggestAction.nextArgs {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg)
	}

	return ActionResult{sb.String()}, []*Consequence{}
}

func MakeAutoSuggestAction(isValidAsIs bool, nextTokens []string) AutoSuggestAction {
	var nextArgs []string

	for _, token := range nextTokens {
		nextArgs = append(nextArgs, token)
	}

	return AutoSuggestAction{isValidAsIs, nextArgs}
}
