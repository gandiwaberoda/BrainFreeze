package commands

import (
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

	parsedFulfillment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parsedFulfillment == nil {
		parsedFulfillment = fulfillments.DefaultHoldFulfillment()
	}

	parsed := HandlingCommand{
		fulfillment: parsedFulfillment,
	}

	return true, &parsed, nil
}

func (i HandlingCommand) GetName() string {
	return "HANDLING"
}

func (i *HandlingCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.EnableHandling()
	i.fulfillment.Tick()
}

func (i HandlingCommand) ShouldClear() bool {
	return i.shouldClear
}

func (i HandlingCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
