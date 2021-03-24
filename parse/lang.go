package parse

import "samvasta.com/bujit/actions"

type ParseContext struct {
	tokens            []string
	currentTokenIndex int
	isValid           bool
	data              map[string]interface{}
}

func (ctx *ParseContext) MoveToNextToken() {
	ctx.currentTokenIndex++
}

func EmptyParseContext(tokens []string) *ParseContext {
	return &ParseContext{tokens, 1, true, make(map[string]interface{})}
}

func (ctx *ParseContext) NextToken() (nextToken string, hasNext bool) {
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

	// Flags
	FROM
	TO
	SET
	PRINT
	FILTER
	ORDER
	BY

	// Models
	ACCOUNT
	ACCOUNT_STATE
	TRANSACTION
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

func ParseExpression(input string) actions.Actioner {
	tokens := Tokenize(input)

	actionTok := tokens[0]

	exact, possible := PossibleMatches(actionTok, ActionTokens)

	if exact == nil {
		return actions.MakeAutoSuggestAction(false, DisplayNames(possible))
	}

	switch exact.Id {
	case NEW:
		return nil
	case LIST:
		return nil
	case DELETE:
		return nil
	case MODIFY:
		return nil
	case DETAIL:
		return nil
	case CLOSE:
		return nil
	case OPEN:
		return nil
	case CONFIGURE:
		return nil
	case HELP:
		return ParseHelpRoot(&HelpContext{EmptyParseContext(tokens), false})
	case EXIT:
		return nil
	case VERSION:
		return nil
	default:
		return nil
	}
}
