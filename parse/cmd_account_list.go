package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/models/output"
)

var listAccountArgs map[int]*TokenPattern = map[int]*TokenPattern{
	ARG_NAME:        MakeOptionalArgToken(ARG_NAME, "n", "name"),
	ARG_DESCRIPTION: MakeOptionalArgToken(ARG_DESCRIPTION, "d", "description"),
	ARG_CATEGORY:    MakeOptionalArgToken(ARG_CATEGORY, "c", "category"),
	ARG_MIN_BALANCE: MakeOptionalArgToken(ARG_MIN_BALANCE, "m", "min-balance"),
	ARG_MAX_BALANCE: MakeOptionalArgToken(ARG_MAX_BALANCE, "x", "max-balance"),
	FLAG_HELP:       makeFlagToken(FLAG_HELP, "h", "help"),
}

type ListAccountContext struct {
	ParseContext
	action                                                             actions_accounts.ListAccountAction
	hasName, hasDescription, hasCategory, hasBalanceMin, hasBalanceMax bool
}

func (ctx ListAccountContext) possibleNextTokens() []*TokenPattern {
	tokens := []*TokenPattern{}

	if !ctx.hasName {
		tokens = append(tokens, listAccountArgs[ARG_NAME])
	}
	if !ctx.hasDescription {
		tokens = append(tokens, listAccountArgs[ARG_DESCRIPTION])
	}
	if !ctx.hasCategory {
		tokens = append(tokens, listAccountArgs[ARG_CATEGORY])
	}
	if !ctx.hasBalanceMin {
		tokens = append(tokens, listAccountArgs[ARG_MIN_BALANCE])
	}
	if !ctx.hasBalanceMax {
		tokens = append(tokens, listAccountArgs[ARG_MAX_BALANCE])
	}

	if !ctx.hasName && !ctx.hasDescription && !ctx.hasCategory && !ctx.hasBalanceMin && !ctx.hasBalanceMax {
		// Only want to accept the help flag if no other args have been seen
		tokens = append(tokens, listAccountArgs[FLAG_HELP])
	}

	return tokens
}

func parseListAccount(context *ListAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
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
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_NAME], ItemNamePattern, "name")
			if suggestion.IsValidAsIs {
				context.action.Name = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_DESCRIPTION:
			context.hasDescription = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_DESCRIPTION], ItemNamePattern, "description")
			if suggestion.IsValidAsIs {
				context.action.Description = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_CATEGORY:
			context.hasCategory = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_CATEGORY], ItemNamePattern, "category")
			if suggestion.IsValidAsIs {
				context.action.CategoryName = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_MIN_BALANCE:
			context.hasBalanceMin = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_MIN_BALANCE], DecimalPattern, "min-balance")
			if suggestion.IsValidAsIs {
				b := moneyValue(value)
				context.action.MinBalance = &b
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_MAX_BALANCE:
			context.hasBalanceMax = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_MAX_BALANCE], DecimalPattern, "max-balance")
			if suggestion.IsValidAsIs {
				b := moneyValue(value)
				context.action.MaxBalance = &b
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case FLAG_HELP:
			return listAccountHelpAction(context)
		}
	} else if context.action.IsValid() {
		return context.action, makeAutoSuggestion(true, nextToken, missingTokens)
	}

	return nil, makeAutoSuggestion(false, "", missingTokens)
}

func listAccountHelpAction(context *ListAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	helpItems := output.EmptyOutputGroup().
		Header("List Accounts Command").
		HorizontalRule("═").
		Header("Description").
		Paragraph("Lists accounts.").
		HorizontalRule("-").
		Header("Syntax: list account [-n=<account-name>] [-c=<category-name>] [-d=<description>] [-m=<min-balance>] [-x=<max-balance>]").
		Indent().
		UnorderedList([]string{
			"name (-n or --name): filter the list of accounts by name. Filters out accounts with names that do not contain the provided value.",
			"category (-c or --category): filter the list of accounts by category. Filters out accounts which do not belong, directly or indirectly, to a category with name containing the provided value.",
			"description (-d or --description): filter the list of accounts by partial description. Filters accounts with descriptions that do not contain the provided value.",
			"min-balance (-m or --min-balance): filter the list of accounts by balance. Filters out accounts with a balance below the provided value.",
			"max-balance (-x or --max-balance): filter the list of accounts by balance. Filters out accounts with a balance above the provided value.",
		}, output.NormalBulletChar).
		Unindent().
		ToSlice()

	return actions.HelpAction{HelpItems: helpItems}, EmptySuggestions
}
