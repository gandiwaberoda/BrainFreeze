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

	return true, &StopCommand{
		fulfillment: parsedFulfilment,
	}, nil
}

func (i StopCommand) GetName() string {
	return "STOP"
}

func (i *StopCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	force.Idle()
}

func (i StopCommand) ShouldClear() bool {
	return false
}

func (i StopCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
