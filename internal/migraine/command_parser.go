package migraine

import (
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

var handlers []func(string, *configuration.FreezeConfig) (bool, commands.CommandInterface) = []func(string, *configuration.FreezeConfig) (bool, commands.CommandInterface){
	commands.ParseIdleCommand,
	commands.ParseWasdCommand,
}

func WhichCommand(intercom models.Intercom, conf *configuration.FreezeConfig) commands.CommandInterface {
	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom.Content, conf)
		if thisIs {
			return cmd
		}
	}

	return nil
}
