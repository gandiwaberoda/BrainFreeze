package migraine

import (
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/pkg/models"
)

var handlers []func(string) (bool, commands.CommandInterface) = []func(string) (bool, commands.CommandInterface){
	commands.ParseIdleCommand,
	commands.ParseWasdCommand,
}

func ParseCommand(intercom models.Intercom) commands.CommandInterface {
	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom.Content)
		if thisIs {
			return cmd
		}
	}

	return nil
}
