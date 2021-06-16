package parser

import (
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

var handlers []func(models.Intercom, *configuration.FreezeConfig) (bool, commands.CommandInterface) = []func(models.Intercom, *configuration.FreezeConfig) (bool, commands.CommandInterface){
	commands.ParseIdleCommand,
	commands.ParseWasdCommand,
}

func WhichCommand(intercom models.Intercom, conf *configuration.FreezeConfig) commands.CommandInterface {
	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom, conf)
		if thisIs {
			return cmd
		}
	}

	return nil
}
