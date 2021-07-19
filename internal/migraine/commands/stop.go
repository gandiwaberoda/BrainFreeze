package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type StopCommand struct {
	fulfillment fulfillments.FulfillmentInterface
}

func DefaultStopCommand() CommandInterface {
	return &IdleCommand{
		fulfillment: fulfillments.DefaultHoldFulfillment(),
	}
}

func ParseStopCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 4 {
	// 	return false, nil
	// }

	// if strings.ToUpper(cmd[:4]) != "STOP" {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "STOP") {
		return false, nil, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultHoldFulfillment()
	}

	return true, &StopCommand{
		fulfillment: parseFulfilment,
	}, nil
}

func (i StopCommand) GetName() string {
	return "STOP"
}

func (i *StopCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.Idle()
	i.fulfillment.Tick()
}

func (i StopCommand) ShouldClear() bool {
	return false
}

func (i StopCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
