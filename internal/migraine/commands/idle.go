package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type IdleCommand struct {
	fulfillment fulfillments.FulfillmentInterface
}

func DefaultIdleCommand() CommandInterface {
	return &IdleCommand{
		fulfillment: fulfillments.DefaultHoldFulfillment(),
	}
}

func ParseIdleCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "IDLE") {
		return false, nil, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultHoldFulfillment()
	}

	return true, &IdleCommand{
		fulfillment: parseFulfilment,
	}, nil
}

func (i IdleCommand) GetName() string {
	return "IDLE"
}

func (i *IdleCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.Idle()
	i.fulfillment.Tick()
}

func (i IdleCommand) ShouldClear() bool {
	return false
}

func (i IdleCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
