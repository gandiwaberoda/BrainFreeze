package commands

import (
	"errors"
	"fmt"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
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

var handlers []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, CommandInterface, error)

func init() {
	handlers = []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, CommandInterface, error){
		ParseIdleCommand,
		ParseWasdCommand,
		ParseLookatCommand,
		ParseHandlingCommand,
		ParseWatchatCommand,
		ParseGetballCommand,
		ParseGotoCommand,
		ParsePlannedCommand,
		ParseReceiveCommand,
		ParseKickCommand,
		ParseApproachCommand,
		ParsePlaylfCommand,
		ParseStopCommand,
		ParsePassingCommand,
	}
}

func WhichCommand(fullbfvid string, conf *configuration.FreezeConfig, state *state.StateAccess) (CommandInterface, error) {
	// splitted := strings.Split(intercom.Content, "/")
	parsed, err := bfvid.ParseCommandSPOK(fullbfvid)
	if err != nil {
		return nil, errors.New(fmt.Sprint("failed to parse command:", err))
	}

	// if len(splitted) < 1 {
	// 	return nil
	// }

	// // CMD berisi commandnya (GETBALL, WATCHAT) dan juga argumen, tanpa RECEIVER, tanpa FULFILLMENT
	// // Misal
	// // all/goto(300,400)/dur(5000)
	// cmd := strings.ToUpper(splitted[0])
	// cmd = strings.ReplaceAll(cmd, "\n", " ")

	for _, isThis := range handlers {
		thisIs, cmd, err := isThis(*parsed, conf, state)
		if thisIs {
			if err != nil {
				return nil, err
			}

			return cmd, nil
		}
	}

	return nil, errors.New(fmt.Sprint("command not found:", parsed.Verb, "\n", fullbfvid))
}
