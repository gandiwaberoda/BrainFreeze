package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
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

func ParseIdleCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 4 {
		return false, nil
	}

	if strings.ToUpper(cmd[:4]) != "IDLE" {
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultHoldFulfillment()
	}

	return true, &IdleCommand{
		fulfillment: parseFulfilment,
	}
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
