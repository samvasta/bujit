package actions

import (
	"samvasta.com/bujit/models"
	"strings"
)

type HelpAction struct {
	helpItems []*models.Text
}

func (helpAction *HelpAction) execute() (ActionResult, []*Consequence) {
	var sb strings.Builder
	return ActionResult{sb.String()}, []*Consequence{}
}

func MakeHelpAction() HelpAction {
	return HelpAction{}
}
