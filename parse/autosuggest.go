package parse

type AutoSuggestion struct {
	IsValidAsIs  bool
	CurrentToken string
	NextArgs     []string
}

func makeAutoSuggestion(isValidAsIs bool, currentToken string, nextTokens []*TokenPattern) AutoSuggestion {
	var nextArgs []string

	for _, token := range nextTokens {
		bestMatch := token.BestMatch(currentToken)
		if bestMatch != "" {
			nextArgs = append(nextArgs, bestMatch)
		}
	}

	return AutoSuggestion{isValidAsIs, currentToken, nextArgs}
}

var EmptySuggestions AutoSuggestion = AutoSuggestion{true, "", []string{}}
