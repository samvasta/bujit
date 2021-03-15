package actions

import (
	"fmt"

	"samvasta.com/bujit/config"
)

type ExitAction struct{}

func (exitAction *ExitAction) execute() (ActionResult, []*Consequence) {
	return ActionResult{"goodbye"}, []*Consequence{}
}

type VersionAction struct{}

func (versionAction *VersionAction) execute() (ActionResult, []*Consequence) {
	return ActionResult{config.Version()}, []*Consequence{}
}

type ConfigureAction struct {
	configureFunc func(interface{})
	value         interface{}
}

func (configureAction *ConfigureAction) execute() (ActionResult, []*Consequence) {
	configureAction.configureFunc(configureAction.value)

	return ActionResult{fmt.Sprintf("Set to %v", configureAction.value)}, []*Consequence{}
}
