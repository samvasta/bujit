package parse

import (
	"samvasta.com/bujit/actions"
	actions_accounts "samvasta.com/bujit/actions/accounts"
	"samvasta.com/bujit/session"
)

type ParseContext struct {
	tokens            []string
	currentTokenIndex int
	isValid           bool
	session           *session.Session
}

func (ctx *ParseContext) moveToNextToken() {
	ctx.currentTokenIndex++
}

func EmptyParseContext(tokens []string, session *session.Session) ParseContext {
	return ParseContext{tokens, 1, true, session}
}

func (ctx *ParseContext) nextToken() (nextToken string, hasNext bool) {
	if ctx.currentTokenIndex >= len(ctx.tokens) {
		return "", false
	}
	return ctx.tokens[ctx.currentTokenIndex], true
}

const (
	// Actions
	EXIT int = iota
	HELP
	VERSION
	CONFIGURE
	NEW
	LIST
	DELETE
	MODIFY
	DETAIL
	CLOSE
	OPEN
	SET
	PRINT

	// Models
	CATEGORY
	ACCOUNT
	ACCOUNT_STATE
	TRANSACTION

	// Args
	ARG_FROM
	ARG_TO
	ARG_NAME
	ARG_DESCRIPTION
	ARG_CATEGORY
	ARG_STARTING_BALANCE
	ARG_MIN_BALANCE
	ARG_MAX_BALANCE

	// Flags
	FLAG_HELP
	FLAG_HARD

	// Misc
	FILTER
	ORDER
	BY
	FROM
	TO
)

var allTokens map[int]*TokenPattern = map[int]*TokenPattern{
	// Commands
	EXIT:      MakeLiteralToken(EXIT, "exit"),
	HELP:      MakeLiteralToken(HELP, "help"),
	VERSION:   MakeLiteralToken(VERSION, "version"),
	CONFIGURE: MakeLiteralToken(CONFIGURE, "config"),
	NEW:       MakeLiteralToken(NEW, "new", "create", "add"),
	LIST:      MakeLiteralToken(LIST, "list", "ls"),
	DELETE:    MakeLiteralToken(DELETE, "delete", "del", "remove", "rm"),
	MODIFY:    MakeLiteralToken(MODIFY, "modify", "mod", "set"),
	DETAIL:    MakeLiteralToken(DETAIL, "detail"),
	CLOSE:     MakeLiteralToken(CLOSE, "close"),
	OPEN:      MakeLiteralToken(OPEN, "open"),

	FROM:  MakeLiteralToken(FROM, "from"),
	TO:    MakeLiteralToken(TO, "to"),
	SET:   MakeLiteralToken(SET, "set"),
	PRINT: MakeLiteralToken(PRINT, "print"),

	FILTER: MakeLiteralToken(FILTER, "filter"),
	ORDER:  MakeLiteralToken(ORDER, "order"),
	BY:     MakeLiteralToken(BY, "by"),

	// Nouns
	CATEGORY:      MakeLiteralToken(CATEGORY, "category", "group"),
	ACCOUNT:       MakeLiteralToken(ACCOUNT, "account", "acct"),
	ACCOUNT_STATE: MakeLiteralToken(ACCOUNT_STATE, "account_state", "acct_state"),
	TRANSACTION:   MakeLiteralToken(TRANSACTION, "transaction", "tran"),
}

var ActionTokens = []*TokenPattern{
	allTokens[NEW],
	allTokens[LIST],
	allTokens[DELETE],
	allTokens[MODIFY],
	allTokens[DETAIL],
	allTokens[CLOSE],
	allTokens[OPEN],
	allTokens[CONFIGURE],
	allTokens[HELP],
	allTokens[EXIT],
	allTokens[VERSION],
}

var ModelTokens = []*TokenPattern{
	allTokens[CATEGORY],
	allTokens[ACCOUNT],
	allTokens[ACCOUNT_STATE],
	allTokens[TRANSACTION],
}

func ParseExpression(input string, session *session.Session) (action actions.Actioner, suggestion AutoSuggestion) {
	tokens := Tokenize(input)

	actionTok := tokens[0]

	exact, possible := PossibleMatches(actionTok, ActionTokens)

	if exact == nil {
		return nil, makeAutoSuggestion(false, actionTok, possible)
	}

	parseContext := EmptyParseContext(tokens, session)

	switch exact.Id {
	case NEW:
		return ParseNew(&parseContext)
	case LIST:
		return ParseList(&parseContext)
	case DELETE:
		return nil, EmptySuggestions
	case MODIFY:
		return nil, EmptySuggestions
	case DETAIL:
		return nil, EmptySuggestions
	case CLOSE:
		return nil, EmptySuggestions
	case OPEN:
		return nil, EmptySuggestions
	case CONFIGURE:
		return nil, EmptySuggestions
	case HELP:
		return parseHelpRoot(&HelpContext{ParseContext: parseContext, verbose: false})
	case EXIT:
		return nil, EmptySuggestions
	case VERSION:
		return nil, EmptySuggestions
	default:
		return nil, EmptySuggestions
	}
}

func ParseNew(context *ParseContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.nextToken()

	if hasNext {
		exact, possible := PossibleMatches(nextToken, ModelTokens)

		if exact == nil {
			return nil, makeAutoSuggestion(false, nextToken, possible)
		}

		switch exact.Id {
		case CATEGORY:
			context.moveToNextToken()
			return nil, EmptySuggestions
		case ACCOUNT:
			context.moveToNextToken()
			return parseNewAccount(
				&NewAccountContext{
					ParseContext: *context,
					action:       actions_accounts.CreateAccountAction{Session: context.session}})
		case ACCOUNT_STATE:
			context.moveToNextToken()
			return nil, EmptySuggestions
		case TRANSACTION:
			context.moveToNextToken()
			return nil, EmptySuggestions
		}
	}

	return nil, makeAutoSuggestion(false, nextToken, ModelTokens)
}

func ParseList(context *ParseContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.nextToken()

	if hasNext {
		exact, possible := PossibleMatches(nextToken, ModelTokens)

		if exact == nil {
			return nil, makeAutoSuggestion(false, nextToken, possible)
		}

		switch exact.Id {
		case CATEGORY:
			context.moveToNextToken()
			return nil, EmptySuggestions
		case ACCOUNT:
			context.moveToNextToken()
			return parseListAccount(
				&ListAccountContext{
					ParseContext: *context,
					action:       actions_accounts.ListAccountAction{Session: context.session}})
		case ACCOUNT_STATE:
			context.moveToNextToken()
			return nil, EmptySuggestions
		case TRANSACTION:
			context.moveToNextToken()
			return nil, EmptySuggestions
		}
	}

	return nil, makeAutoSuggestion(false, nextToken, ModelTokens)
}
