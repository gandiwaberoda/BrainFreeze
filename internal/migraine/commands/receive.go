package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type ReceiveCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseReceiveCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 7 {
	// 	return false, nil
	// }

	// if !strings.EqualFold(cmd[:7], "RECEIVE") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "RECEIVE") {
		return false, nil, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultGotballFulfillment(curstate)
	}
	parsed := ReceiveCommand{
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed, nil
}

func (i ReceiveCommand) GetName() string {
	return "RECEIVE"
}

func (i *ReceiveCommand) Tick(force *models.Force, curstate *state.StateAccess) {
	// _, target := state.GetTransformByKey(i.Target)
	var target models.Transform

	if !curstate.GetState().BallTransformExpired {
		// Utamakan ngelihat ke bola
		target = curstate.GetState().BallTransform
	} else {
		// Jika bolanya LOST, lihat ke temen
		if i.conf.Robot.Color == configuration.CYAN {
			target = curstate.GetState().MagentaTransform
		} else if i.conf.Robot.Color == configuration.MAGENTA {
			target = curstate.GetState().CyanTransform
		} else {
			panic(i.conf.Robot.Color + " is not a valid robot color")
		}
	}

	TockLookat(target, *i.conf, force, curstate)
	force.EnableHandling()

	register := state.NewRegister()
	if target.RobROT <= models.Degree(i.conf.CommandParameter.LookatToleranceDeg) {
		register.ReadyReceive = 1.0
	} else {
		register.ReadyReceive = 0.0
	}
	curstate.UpdateRegisterState(register)

	i.fulfillment.Tick()
}

func (i ReceiveCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
