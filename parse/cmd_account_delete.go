package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models/output"
)

var deleteAccountArgs map[int]*TokenPattern = map[int]*TokenPattern{
	ARG_NAME:  MakeOptionalArgToken(ARG_NAME, "n", "name"),
	FLAG_HARD: makeFlagToken(FLAG_HARD, "d", "hard"),
	FLAG_HELP: makeFlagToken(FLAG_HELP, "h", "help"),
}

type DeleteAccountContext struct {
	ParseContext
	action             actions_accounts.DeleteAccountAction
	hasName, hasIsHard bool
}

func (ctx DeleteAccountContext) possibleNextTokens() []*TokenPattern {
	tokens := []*TokenPattern{}

	if !ctx.hasName {
		tokens = append(tokens, deleteAccountArgs[ARG_NAME])

		// Only want to accept the help flag if no other args have been seen
		tokens = append(tokens, deleteAccountArgs[FLAG_HELP])
	} else if !ctx.hasIsHard {
		tokens = append(tokens, deleteAccountArgs[FLAG_HARD])
	}

	return tokens
}

func parseDeleteAccount(context *DeleteAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.nextToken()

	missingTokens := []*TokenPattern{}

	if !context.hasName {
		missingTokens = append(missingTokens, deleteAccountArgs[ARG_NAME])
		missingTokens = append(missingTokens, deleteAccountArgs[FLAG_HELP])
	} else if !context.hasIsHard {
		missingTokens = append(missingTokens, deleteAccountArgs[FLAG_HARD])
	}

	if hasNext {
		possibleTokens := context.possibleNextTokens()
		exact, possible := PossibleMatches(nextToken, possibleTokens)

		if exact == nil {
			return nil, makeAutoSuggestion(false, nextToken, possible)
		}

		switch exact.Id {
		case ARG_NAME:
			context.hasName = true
			value, suggestion := parseOptionalArg(&context.ParseContext, deleteAccountArgs[ARG_NAME], ItemNamePattern, "name")
			if suggestion.IsValidAsIs {
				context.action.Name = itemNameValue(value)
				return parseDeleteAccount(context)
			} else {
				return nil, suggestion
			}
		case FLAG_HARD:
			context.hasIsHard = true
			context.action.IsHardDelete = true
			//No more possible tokens. May as well return the action now
			return context.action, makeAutoSuggestion(true, nextToken, missingTokens)
		case FLAG_HELP:
			return deleteAccountHelpAction(context)
		}
	} else if context.action.IsValid() {
		return context.action, makeAutoSuggestion(true, nextToken, missingTokens)
	}

	return nil, makeAutoSuggestion(false, "", missingTokens)
}

func deleteAccountHelpAction(context *DeleteAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	helpItems := output.EmptyOutputGroup().
		Header("Delete Account Command").
		HorizontalRule("═").
		Header("Description").
		Paragraph("Deletes an account by name. By default, this command is a 'soft delete' and can be undone later with the 'open account' command.").
		HorizontalRule("-").
		EmptyLines(1).
		Header("Syntax: delete account <name> [-d or --hard]").
		Indent().
		Paragraph("hard (-d or --hard): permanently deletes all account data. This cannot be undone.").
		Unindent().
		ToSlice()

	return actions.HelpAction{HelpItems: helpItems}, EmptySuggestions
}
