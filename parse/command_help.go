package parse

import (
	"fmt"

	"samvasta.com/bujit/actions"
	"samvasta.com/bujit/models/output"
)

type HelpContext struct {
	ParseContext
	verbose bool
}

var VerboseToken *TokenPattern = MakeFlagToken(1, "v", "verbose")
var parseRootNextTokens = []*TokenPattern{VerboseToken}

func ParseHelpRoot(context *HelpContext) actions.Actioner {
	nextToken, hasNext := context.NextToken()
	fmt.Println("next: " + nextToken)

	if hasNext {
		exact, possible := PossibleMatches(nextToken, parseRootNextTokens)

		if exact == nil {
			fmt.Println("exact is nil")
			suggestions := DisplayNames(possible)
			fmt.Printf("suggestions1 %s\n", suggestions)
			return actions.MakeAutoSuggestAction(false, suggestions)
		}

		switch exact.Id {
		case VerboseToken.Id:
			context.verbose = true
			context.MoveToNextToken()
			return ParseHelpRoot(context)
		default:
			return actions.MakeAutoSuggestAction(false, DisplayNames(possible))
		}
	} else if VerboseToken.Matches(nextToken) {
		// return verbose help
		return &actions.HelpAction{HelpItems: []output.Helper{}}
	} else {
		// return terse help
		helpItems := output.EmptyOutputGroup().
			Header("Bujit General Help").
			HorizontalRule("═").
			Paragraph("Available Commands").
			UnorderedList(DisplayNames(ActionTokens), output.NormalBulletChar).
			ToSlice()

		return &actions.HelpAction{HelpItems: helpItems}
	}

}
