package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type KickCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseKickCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 4 || !strings.EqualFold(cmd[:4], "KICK") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "KICK") {
		return false, nil, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultLostballFulfillment(curstate)
	}
	parsed := KickCommand{
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed, nil
}

func (i KickCommand) GetName() string {
	return "KICK"
}

func (i *KickCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.Kick()
	force.EnableHandling()
	i.fulfillment.Tick()
}

func (i KickCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
