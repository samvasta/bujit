package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models/output"
)

var newAccountArgs map[int]*TokenPattern = map[int]*TokenPattern{
	ARG_NAME:             MakeArgToken(ARG_NAME, "name", ItemNamePattern),
	ARG_DESCRIPTION:      MakeOptionalArgToken(ARG_DESCRIPTION, "d", "description"),
	ARG_CATEGORY:         MakeOptionalArgToken(ARG_CATEGORY, "c", "category"),
	ARG_STARTING_BALANCE: MakeOptionalArgToken(ARG_STARTING_BALANCE, "b", "balance"),
	FLAG_HELP:            makeFlagToken(FLAG_HELP, "h", "help"),
}

type NewAccountContext struct {
	ParseContext
	action                                                   actions_accounts.CreateAccountAction
	hasName, hasDescription, hasCategory, hasStartingBalance bool
}

func (ctx NewAccountContext) possibleNextTokens() []*TokenPattern {
	tokens := []*TokenPattern{}
	if !ctx.hasName {
		tokens = append(tokens, newAccountArgs[ARG_NAME])

		// Only want to accept the help flag if no other args have been seen
		tokens = append(tokens, newAccountArgs[FLAG_HELP])
	} else {
		if !ctx.hasDescription {
			tokens = append(tokens, newAccountArgs[ARG_DESCRIPTION])
		}
		if !ctx.hasCategory {
			tokens = append(tokens, newAccountArgs[ARG_CATEGORY])
		}
		if !ctx.hasStartingBalance {
			tokens = append(tokens, newAccountArgs[ARG_STARTING_BALANCE])
		}
	}

	return tokens
}

func parseNewAccount(context *NewAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.nextToken()

	missingTokens := context.possibleNextTokens()

	if hasNext {
		possibleTokens := missingTokens
		exact, possible := PossibleMatches(nextToken, possibleTokens)

		if exact == nil {
			return nil, makeAutoSuggestion(false, nextToken, possible)
		}

		switch exact.Id {
		case ARG_NAME:
			context.hasName = true
			context.action.Name = itemNameValue(nextToken)
			context.moveToNextToken()
			return parseNewAccount(context)
		case ARG_DESCRIPTION:
			context.hasDescription = true
			value, suggestion := parseOptionalArg(&context.ParseContext, newAccountArgs[ARG_DESCRIPTION], ItemNamePattern, "description")
			if suggestion.IsValidAsIs {
				context.action.Description = itemNameValue(value)
				return parseNewAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_CATEGORY:
			context.hasCategory = true
			value, suggestion := parseOptionalArg(&context.ParseContext, newAccountArgs[ARG_CATEGORY], ItemNamePattern, "category")
			if suggestion.IsValidAsIs {
				context.action.CategoryName = itemNameValue(value)
				return parseNewAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_STARTING_BALANCE:
			context.hasStartingBalance = true
			value, suggestion := parseOptionalArg(&context.ParseContext, newAccountArgs[ARG_STARTING_BALANCE], DecimalPattern, "balance")
			if suggestion.IsValidAsIs {
				context.action.StartingBalance = moneyValue(value)
				return parseNewAccount(context)
			} else {
				return nil, suggestion
			}
		case FLAG_HELP:
			return newAccountHelpAction(context)
		}
	} else if context.action.IsValid() {
		return context.action, makeAutoSuggestion(true, nextToken, missingTokens)
	}

	return nil, makeAutoSuggestion(false, "", missingTokens)
}

func newAccountHelpAction(context *NewAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	helpItems := output.EmptyOutputGroup().
		Header("Create New Account Command").
		HorizontalRule("═").
		Header("Description").
		Paragraph("Create a new account.").
		HorizontalRule("-").
		Header("Syntax: new account <name> [-c=<category-name>] [-d=<description>] [-b=<starting-balance>]").
		ToSlice()

	return actions.HelpAction{HelpItems: helpItems}, EmptySuggestions
}
