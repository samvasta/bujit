package parse

import (
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models/output"
)

type HelpContext struct {
	ParseContext
	verbose bool
}

var VerboseToken *TokenPattern = MakeFlagToken(1, "v", "verbose")
var parseRootNextTokens = []*TokenPattern{VerboseToken}

func ParseHelpRoot(context *HelpContext) (action actions.Actioner, suggestion AutoSuggestion) {
	nextToken, hasNext := context.NextToken()

	if hasNext {
		exact, possible := PossibleMatches(nextToken, parseRootNextTokens)

		if exact == nil {
			return nil, MakeAutoSuggestion(false, DisplayNames(possible))
		}

		switch exact.Id {
		case VerboseToken.Id:
			context.verbose = true
			context.MoveToNextToken()
			return ParseHelpRoot(context)
		default:
			return nil, MakeAutoSuggestion(false, DisplayNames(possible))
		}
	} else if VerboseToken.Matches(nextToken) {
		// return verbose help
		return actions.HelpAction{HelpItems: []output.Helper{}}, EmptySuggestions
	} else {
		// return terse help
		helpItems := output.EmptyOutputGroup().
			Header("Bujit General Help").
			HorizontalRule("═").
			Paragraph("Available Commands").
			UnorderedList(DisplayNames(ActionTokens), output.NormalBulletChar).
			ToSlice()

		return actions.HelpAction{HelpItems: helpItems}, EmptySuggestions
	}

}
