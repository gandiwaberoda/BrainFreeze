package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type HandlingCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

func ParseHandlingCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig) (bool, CommandInterface) {
	if len(cmd) < 8 {
		return false, nil
	}

	if strings.ToUpper(cmd[:8]) != "HANDLING" {
		return false, &HandlingCommand{}
	}

	parsed := HandlingCommand{}
	parsedFulfillment := fulfillments.WhichFulfillment(intercom, conf)
	if parsedFulfillment == nil {
		parsedFulfillment = fulfillments.DefaultHoldFulfillment()
	}
	parsed.fulfillment = parsedFulfillment

	return true, &parsed
}

func (i HandlingCommand) GetName() string {
	return "HANDLING"
}

func (i *HandlingCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.EnableHandling()
	i.shouldClear = i.fulfillment.Tick(state)
}

func (i HandlingCommand) ShouldClear() bool {
	return i.shouldClear
}

func (i HandlingCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
