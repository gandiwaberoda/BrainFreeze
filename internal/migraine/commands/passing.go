package commands

// TODO: Tunggu temen posisinya ready receive, lalu kick, cari posisi juga maybe
// TODO: Bolehkan ngoper ke robot lain, tambahkan parameter Target

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PassingCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParsePassingCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 7 || !strings.EqualFold(cmd[:7], "PASSING") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "PASSING") {
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

	parsed := PassingCommand{
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i PassingCommand) GetName() string {
	return "PASSING"
}

func (i *PassingCommand) Tick(force *models.Force, curstate *state.StateAccess) {
	i.fulfillment.Tick()

	var target models.Transform
	var targetColor string
	if i.conf.Robot.Color == configuration.CYAN {
		target = curstate.GetState().MagentaTransform
		targetColor = "MAGENTA"
	} else if i.conf.Robot.Color == configuration.MAGENTA {
		target = curstate.GetState().CyanTransform
		targetColor = "CYAN"
	} else {
		panic(fmt.Sprint(i.conf.Robot.Color, "is not a valid color"))
	}

	TockLookat(target, *i.conf, force, curstate)
	if math.Abs(float64(target.RobROT)) > float64(i.conf.CommandParameter.LookatToleranceDeg) {
		return
	}

	readyNerima, err := curstate.GetOtherRegisterByIdentifier(targetColor, state.READY_RECEIVED)
	_ = readyNerima
	if err != nil {
		fmt.Println("failed getting register by identifier:", err)
		return
	}

	if readyNerima == 1.0 {
		force.Kick()
		force.EnableHandling()
	}

}

func (i PassingCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
