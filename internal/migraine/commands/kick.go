package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type KickCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseKickCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 4 || !strings.EqualFold(cmd[:4], "KICK") {
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultLostballFulfillment(curstate)
	}
	parsed := KickCommand{
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed
}

func (i KickCommand) GetName() string {
	return "KICK"
}

func (i *KickCommand) Tick(force *models.Force, state *state.StateAccess) {
	force.Kick()
	i.fulfillment.Tick()
}

func (i KickCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
