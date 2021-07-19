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

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultGotballFulfillment(curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := ReceiveCommand{
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i ReceiveCommand) GetName() string {
	return "RECEIVE"
}

func (i *ReceiveCommand) Tick(force *models.Force, curstate *state.StateAccess) {
	// _, target := state.GetTransformByKey(i.Target)

	var target models.Transform
	var partner models.Transform

	// Jika bolanya LOST, lihat ke temen
	if i.conf.Robot.Color == configuration.CYAN {
		partner = curstate.GetState().MagentaTransform
	} else if i.conf.Robot.Color == configuration.MAGENTA {
		partner = curstate.GetState().CyanTransform
	} else {
		panic(i.conf.Robot.Color + " is not a valid robot color")
	}

	if !curstate.GetState().BallTransformExpired {
		// Utamakan ngelihat ke bola
		target = curstate.GetState().BallTransform
	} else {
		target = partner
	}

	TockLookat(target, *i.conf, force, curstate)
	force.EnableHandling()

	register := state.NewRegister()
	// Register ReadyReceive berfungsi untuk mengatakan kepada partner bahwa aku siap nerima bola
	// So, targetnya bukan bola, tapi sedang melihat partner yang akan ngoper
	if partner.RobROT <= models.Degree(i.conf.CommandParameter.LookatToleranceDeg) {
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
