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

type HandlingCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

func ParseHandlingCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 8 {
	// 	return false, nil
	// }

	// if strings.ToUpper(cmd[:8]) != "HANDLING" {
	// 	return false, &HandlingCommand{}
	// }
	if !strings.EqualFold(cmd.Verb, "HANDLING") {
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

	parsed := HandlingCommand{
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i HandlingCommand) GetName() string {
	return "HANDLING"
}

func (i *HandlingCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	force.EnableHandling()
}

func (i HandlingCommand) ShouldClear() bool {
	return i.shouldClear
}

func (i HandlingCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
