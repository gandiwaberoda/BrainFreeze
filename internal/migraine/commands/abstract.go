package commands

import (
	"fmt"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type CommandInterface interface {
	GetName() string
	Tick(*models.Force, *state.StateAccess)
	// ShouldClear() bool
	GetFulfillment() fulfillments.FulfillmentInterface
}

var handlers []func(models.Intercom, string, *configuration.FreezeConfig, *state.StateAccess) (bool, CommandInterface) = []func(models.Intercom, string, *configuration.FreezeConfig, *state.StateAccess) (bool, CommandInterface){
	ParseIdleCommand,
	ParseWasdCommand,
	ParseLookatCommand,
	ParseHandlingCommand,
	ParseWatchatCommand,
	ParseGetballCommand,
	ParseGotoCommand,
	ParsePlannedCommand,
	ParseReceiveCommand,
}

func WhichCommand(intercom models.Intercom, conf *configuration.FreezeConfig, state *state.StateAccess) CommandInterface {
	splitted := strings.Split(intercom.Content, "/")

	if len(splitted) < 1 {
		return nil
	}

	// CMD berisi commandnya (GETBALL, WATCHAT) dan juga argumen, tanpa RECEIVER, tanpa FULFILLMENT
	// Misal
	// all/goto(300,400)/dur(5000)
	cmd := strings.ToUpper(splitted[0])
	cmd = strings.ReplaceAll(cmd, "\n", " ")

	for _, isThis := range handlers {
		thisIs, cmd := isThis(intercom, cmd, conf, state)
		if thisIs {
			return cmd
		}
	}

	fmt.Println("Gak ketemu:", cmd)
	return nil
}
