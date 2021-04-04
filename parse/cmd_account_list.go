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

	if len(tokens) == 0 {
		// Only want to accept the help flag if no other args have been seen
		tokens = append(tokens, listAccountArgs[FLAG_HELP])
	}

	return tokens
}

func parseListAccount(context *ListAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.nextToken()

	missingTokens := []*TokenPattern{}

	if !context.hasName {
		missingTokens = append(missingTokens, listAccountArgs[ARG_NAME])
	}
	if !context.hasDescription {
		missingTokens = append(missingTokens, listAccountArgs[ARG_DESCRIPTION])
	}

	if !context.hasCategory {
		missingTokens = append(missingTokens, listAccountArgs[ARG_CATEGORY])
	}

	if !context.hasBalanceMin {
		missingTokens = append(missingTokens, listAccountArgs[ARG_MIN_BALANCE])
	}
	if !context.hasBalanceMax {
		missingTokens = append(missingTokens, listAccountArgs[ARG_MAX_BALANCE])
	}

	if !context.hasName && !context.hasDescription && !context.hasCategory && !context.hasBalanceMin && !context.hasBalanceMax {
		missingTokens = append(missingTokens, listAccountArgs[FLAG_HELP])
	}

	if hasNext {
		possibleTokens := context.possibleNextTokens()
		exact, possible := PossibleMatches(nextToken, possibleTokens)

		if exact == nil {
			return nil, makeAutoSuggestion(false, DisplayNames(possible))
		}

		switch exact.Id {
		case ARG_NAME:
			context.hasName = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_NAME], ItemNamePattern, "name")
			if suggestion.isValidAsIs {
				context.action.Name = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_DESCRIPTION:
			context.hasDescription = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_DESCRIPTION], ItemNamePattern, "description")
			if suggestion.isValidAsIs {
				context.action.Description = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_CATEGORY:
			context.hasCategory = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_CATEGORY], ItemNamePattern, "category")
			if suggestion.isValidAsIs {
				context.action.CategoryName = itemNameValue(value)
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_MIN_BALANCE:
			context.hasBalanceMin = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_MIN_BALANCE], DecimalPattern, "min-balance")
			if suggestion.isValidAsIs {
				b := moneyValue(value)
				context.action.MinBalance = &b
				return parseListAccount(context)
			} else {
				return nil, suggestion
			}
		case ARG_MAX_BALANCE:
			context.hasBalanceMax = true
			value, suggestion := parseOptionalArg(&context.ParseContext, listAccountArgs[ARG_MAX_BALANCE], DecimalPattern, "max-balance")
			if suggestion.isValidAsIs {
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
		return context.action, makeAutoSuggestion(true, DisplayNames(missingTokens))
	}

	return nil, makeAutoSuggestion(false, DisplayNames(missingTokens))
}

func listAccountHelpAction(context *ListAccountContext) (action actions.Actioner, suggestion AutoSuggestion) {
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
