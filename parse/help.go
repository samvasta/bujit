package parse

import (
	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models/output"
)

type HelpContext struct {
	parseContext *ParseContext
	verbose      bool
}

var VerboseToken *TokenPattern = MakeFlagToken(1, "v", "verbose")

func ParseHelpRoot(context *HelpContext) actions.Actioner {
	token, isValid := context.parseContext.NextToken()

	if isValid {
		exact, _ := PossibleMatches(token, []*TokenPattern{VerboseToken})

		switch exact.Id {
		case VerboseToken.Id:
			context.verbose = true
			context.parseContext.MoveToNextToken()
			return ParseHelpRoot(context)
		default:
			return actions.MakeAutoSuggestAction(false, []string{VerboseToken.DisplayName})
		}
	} else if VerboseToken.Matches(token) {
		// return verbose help
		return &actions.HelpAction{HelpItems: []output.Helper{}}
	} else {
		// return terse help
		helpItems := output.EmptyOutputGroup().
			Header("Bujit General Help").
			HorizontalRule("═").
			Paragraph("Available Commands").
			UnorderedList(DisplayNames(ActionTokens)).
			ToSlice()

		return &actions.HelpAction{HelpItems: helpItems}
	}

}
