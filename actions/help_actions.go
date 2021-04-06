package actions

import (
	"samvasta.com/bujit/models/output"
)

type HelpAction struct {
	HelpItems []output.Helper
}

func (helpAction HelpAction) Execute() (ActionResult, []*Consequence) {

	return ActionResult{Output: helpAction.HelpItems, IsSuccessful: true}, []*Consequence{}
}
