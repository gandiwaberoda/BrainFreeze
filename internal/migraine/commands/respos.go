// Untuk reset ke starting position

package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type ResposCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseResposCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "RESPOS") {
		return false, nil, nil
	}

	parsed := ResposCommand{
		conf:        conf,
		fulfillment: fulfillments.DefaultDurationFulfillment(),
	}

	return true, &parsed, nil
}

func (i ResposCommand) GetName() string {
	return "RESPOS"
}

func (i *ResposCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	force.DoReset()
}

func (i ResposCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
