package commands

import (
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type CommandInterface interface {
	GetName() string
	Tick(*models.Force, *state.StateAccess)
	ShouldClear() bool
}

var handlers []func(models.Intercom, *configuration.FreezeConfig) (bool, CommandInterface) = []func(models.Intercom, *configuration.FreezeConfig) (bool, CommandInterface){
	ParseIdleCommand,
	ParseWasdCommand,
}

func WhichCommand(intercom models.Intercom, conf *configuration.FreezeConfig) CommandInterface {
	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom, conf)
		if thisIs {
			return cmd
		}
	}

	return nil
}
