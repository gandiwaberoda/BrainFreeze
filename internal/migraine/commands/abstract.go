package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type CommandInterface interface {
	GetName() string
	Tick(*models.Force, *state.StateAccess)
	ShouldClear() bool
	GetFulfillment() fulfillments.FulfillmentInterface
}

var handlers []func(models.Intercom, string, *configuration.FreezeConfig) (bool, CommandInterface) = []func(models.Intercom, string, *configuration.FreezeConfig) (bool, CommandInterface){
	ParseIdleCommand,
	ParseWasdCommand,
	ParseLookatCommand,
	ParseHandlingCommand,
	ParseWatchatCommand,
}

func WhichCommand(intercom models.Intercom, conf *configuration.FreezeConfig) CommandInterface {
	splitted := strings.Split(intercom.Content, "/")

	if len(splitted) < 1 {
		return nil
	}

	cmd := strings.ToUpper(splitted[0])

	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom, cmd, conf)
		if thisIs {
			return cmd
		}
	}

	return nil
}
