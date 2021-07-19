package commands

import (
	"errors"
	"fmt"
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

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultHoldFulfillment()
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	return true, &IdleCommand{
		fulfillment: parsedFulfilment,
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
