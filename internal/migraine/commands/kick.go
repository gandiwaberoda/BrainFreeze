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

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultLostballFulfillment(curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := KickCommand{
		conf:        conf,
		fulfillment: parsedFulfilment,
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
