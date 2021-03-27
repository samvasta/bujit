package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
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
	}
	if !ctx.hasDescription {
		tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_DESCRIPTION])
	}
	if !ctx.hasCategory {
		tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_CATEGORY])
	}
	if !ctx.hasStartingBalance {
		tokens = append(tokens, newAccountArgs[ACCOUNT_ARG_STARTING_BALANCE])
	}

	// Only want to accept the help flag if no other args have been seen
	if len(tokens) == 0 {
		tokens = append(tokens, newAccountArgs[ACCOUNT_FLAG_HELP])
	}

	return tokens
}

func ParseNewAccount(context *NewAccountContext) actions.Actioner {
	nextToken, hasNext := context.NextToken()

	if hasNext {
		possibleTokens := context.PossibleNextTokens()
		exact, _ := PossibleMatches(nextToken, possibleTokens)

		switch exact.Id {
		case ACCOUNT_ARG_NAME:
			context.hasName = true
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_DESCRIPTION:
			context.hasDescription = true
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_CATEGORY:
			context.hasCategory = true
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_ARG_STARTING_BALANCE:
			context.hasStartingBalance = true
			context.MoveToNextToken()
			return ParseNewAccount(context)
		case ACCOUNT_FLAG_HELP:
			// TODO: Create help action specific to the "create account" command
			return nil
		}
	} else if context.action.IsValid() {
		return &context.action
	}
	missingTokens := []*TokenPattern{}

	if context.action.Name == "" {
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_NAME])
	}

	if context.action.Description == "" {
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_DESCRIPTION])
	}

	if context.action.CategoryName == "" {
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_CATEGORY])
	}

	if context.action.StartingBalance.Value() == 0 {
		missingTokens = append(missingTokens, newAccountArgs[ACCOUNT_ARG_STARTING_BALANCE])
	}

	return actions.MakeAutoSuggestAction(false, DisplayNames(missingTokens))
}
