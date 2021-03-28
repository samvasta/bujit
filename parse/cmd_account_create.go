package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models/output"
)

const (
	ACCOUNT_ARG_NAME = iota
	ACCOUNT_ARG_DESCRIPTION
	ACCOUNT_ARG_CATEGORY
	ACCOUNT_ARG_STARTING_BALANCE
	ACCOUNT_FLAG_HELP
)

var newAccountArgs map[int]*TokenPattern = map[int]*TokenPattern{
	ACCOUNT_ARG_NAME:             MakeArgToken(ACCOUNT_ARG_NAME, "name", ItemNamePattern),
	ACCOUNT_ARG_DESCRIPTION:      MakeOptionalArgToken(ACCOUNT_ARG_NAME, "d", "description", "str", ItemNamePattern),
	ACCOUNT_ARG_CATEGORY:         MakeOptionalArgToken(ACCOUNT_ARG_NAME, "c", "category", "str", ItemNamePattern),
	ACCOUNT_ARG_STARTING_BALANCE: MakeOptionalArgToken(ACCOUNT_ARG_NAME, "b", "balance", "amount", DecimalPattern),
	ACCOUNT_FLAG_HELP:            MakeFlagToken(ACCOUNT_FLAG_HELP, "h", "help"),
}

type NewAccountContext struct {
	ParseContext
	action                                                   actions_accounts.CreateAccountAction
	hasName, hasDescription, hasCategory, hasStartingBalance bool
}

func (ctx NewAccountContext) PossibleNextTokens() []*TokenPattern {
	tokens := []*TokenPattern{}
	if !ctx.hasName {
		tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_NAME])

		// Only want to accept the help flag if no other args have been seen
		tokens = append(tokens, newAccountArgs[ACCOUNT_FLAG_HELP])
	} else {
		if !ctx.hasDescription {
			tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_DESCRIPTION])
		}
		if !ctx.hasCategory {
			tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_CATEGORY])
		}
		if !ctx.hasStartingBalance {
			tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_STARTING_BALANCE])
		}
	}

	return tokens
}

func ParseNewAccount(context *NewAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.NextToken()

	missingTokens := []*TokenPattern{}

	if context.action.Name == "" {
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_NAME])
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_FLAG_HELP])
	} else {
		if context.action.Description == "" {
			missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_DESCRIPTION])
		}

		if context.action.CategoryName == "" {
			missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_CATEGORY])
		}

		if context.action.StartingBalance.Value() == 0 {
			missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_STARTING_BALANCE])
		}
	}

	if hasNext {
		possibleTokens := context.PossibleNextTokens()
		exact, possible := PossibleMatches(nextToken, possibleTokens)

		if exact == nil {
			return nil, MakeAutoSuggestion(false, DisplayNames(possible))
		}

		switch exact.Id {
		case ACCOUNT_ARG_NAME:
			context.hasName = true
			context.action.Name = ItemNameValue(nextToken)
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_DESCRIPTION:
			context.hasDescription = true
			context.action.Description = ItemNameValue(nextToken)
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_CATEGORY:
			context.hasCategory = true
			context.action.CategoryName = ItemNameValue(nextToken)
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_STARTING_BALANCE:
			context.hasStartingBalance = true
			context.action.StartingBalance = MoneyValue(nextToken)
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_FLAG_HELP:
			return NewAccountHelpAction(context)
		}
	} else if context.action.IsValid() {
		return context.action, MakeAutoSuggestion(true, DisplayNames(missingTokens))
	}

	return nil, MakeAutoSuggestion(false, DisplayNames(missingTokens))
}

func NewAccountHelpAction(context *NewAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
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
