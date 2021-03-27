package actions

import (
	"encoding/json"

	"samvasta.com/bujit/models/output"
)

type HelpAction struct {
	HelpItems []output.Helper
}

func (helpAction *HelpAction) Execute() (ActionResult, []*Consequence) {
	jsonHelpItems, _ := json.Marshal(helpAction.HelpItems)

	return ActionResult{Output: string(jsonHelpItems), Suggestions: []string{}}, []*Consequence{}
}
