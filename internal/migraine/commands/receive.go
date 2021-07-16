package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type ReceiveCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseReceiveCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 7 {
		return false, nil
	}

	if !strings.EqualFold(cmd[:7], "RECEIVE") {
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultGotballFulfillment(curstate)
	}
	parsed := ReceiveCommand{
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed
}

func (i ReceiveCommand) GetName() string {
	return "RECEIVE"
}

func (i *ReceiveCommand) Tick(force *models.Force, state *state.StateAccess) {
	// _, target := state.GetTransformByKey(i.Target)
	var target models.Transform

	if !state.GetState().BallTransformExpired {
		// Utamakan ngelihat ke bola
		target = state.GetState().BallTransform
	} else {
		// Jika bolanya LOST, lihat ke temen
		if i.conf.Robot.Color == configuration.CYAN {
			target = state.GetState().MagentaTransform
		} else if i.conf.Robot.Color == configuration.MAGENTA {
			target = state.GetState().CyanTransform
		} else {
			panic(i.conf.Robot.Color + " is not a valid robot color")
		}
	}

	TockLookat(target, *i.conf, force, state)
	force.EnableHandling()

	i.fulfillment.Tick()
}

func (i ReceiveCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
