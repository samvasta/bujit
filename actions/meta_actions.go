package actions

import (
	"fmt"

	"samvasta.com/bujit/config"
)

type ExitAction struct{}

func (exitAction ExitAction) Execute() (ActionResult, []*Consequence) {
	return ActionResult{"goodbye", true}, []*Consequence{}
}

type VersionAction struct{}

func (versionAction VersionAction) Execute() (ActionResult, []*Consequence) {
	return ActionResult{config.Version(), true}, []*Consequence{}
}

type ConfigureAction struct {
	configureFunc func(interface{})
	value         interface{}
}

func (configureAction ConfigureAction) Execute() (ActionResult, []*Consequence) {
	configureAction.configureFunc(configureAction.value)

	return ActionResult{fmt.Sprintf("Set to %v", configureAction.value), true}, []*Consequence{}
}
