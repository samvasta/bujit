package parse

import "samvasta.com/bujit/actions"

func ParseHelpRoot(context *ParseContext) actions.Actioner {
	token := context.NextToken()

	exact, possible := PossibleMatches(token, ActionTokens)

	if exact == nil {
		return actions.MakeAutoSuggestAction(token == "", DisplayNames(possible))
	}

	switch exact.Id {
	default:
		return nil
	}
}
