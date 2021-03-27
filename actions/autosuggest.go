package actions

type AutoSuggestAction struct {
	isValidAsIs bool
	nextArgs    []string
}

func (autoSuggestAction AutoSuggestAction) Execute() (result ActionResult, consequences []*Consequence) {
	return ActionResult{"", autoSuggestAction.nextArgs}, []*Consequence{}
}

func MakeAutoSuggestAction(isValidAsIs bool, nextTokens []string) AutoSuggestAction {
	var nextArgs []string

	for _, token := range nextTokens {
		nextArgs = append(nextArgs, token)
	}

	return AutoSuggestAction{isValidAsIs, nextArgs}
}
