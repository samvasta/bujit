package parse

type AutoSuggestion struct {
	isValidAsIs bool
	nextArgs    []string
}

func MakeAutoSuggestion(isValidAsIs bool, nextTokens []string) AutoSuggestion {
	var nextArgs []string

	for _, token := range nextTokens {
		nextArgs = append(nextArgs, token)
	}

	return AutoSuggestion{isValidAsIs, nextArgs}
}

var EmptySuggestions AutoSuggestion = AutoSuggestion{true, []string{}}
