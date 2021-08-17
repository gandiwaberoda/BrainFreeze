package commands

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type LookatCommand struct {
	Target      string
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseLookatCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 6 {
	// 	return false, nil
	// }

	// if len(cmd) < 6 || !strings.EqualFold(cmd[:6], "LOOKAT") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "LOOKAT") {
		return false, nil, nil
	}

	// re, _ := regexp.Compile(`\(([A-Za-z0-9]+)\)`)
	// foundParam := re.FindString(cmd)
	// foundParam = strings.ReplaceAll(foundParam, "(", "")
	// foundParam = strings.ReplaceAll(foundParam, ")", "")

	// fmt.Println("zzz", foundParam)

	target := "BALL"

	ms := 0
	if len(cmd.Parameter) == 1 && cmd.Parameter[0] != "" {
		target = cmd.Parameter[0]
	} else if len(cmd.Parameter) == 2 {
		target = cmd.Parameter[0]
		ms_, err := strconv.Atoi(cmd.Parameter[1])
		if err != nil {
			return true, nil, errors.New("lookat command can't parse int of ms to clear")
		}
		ms = ms_
	} else if len(cmd.Parameter) > 2 {
		return true, nil, errors.New("lookat command require either NONE, 1 or 2 parameter")
	}

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		return true, nil, errors.New("lookat target key is not acceptable")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultGlancedFulfillment(target, ms, curstate, conf)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := LookatCommand{
		Target:      target,
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i LookatCommand) GetName() string {
	return "LOOKAT:" + string(i.Target)
}

func TockLookat(target models.Transform, conf configuration.FreezeConfig, force *models.Force, state *state.StateAccess) {
	rotForce := target.RobROT

	if math.Abs(float64(rotForce)) < float64(conf.Mecha.RotationForceMinRange) {
		if rotForce < 0 {
			rotForce = models.Degree(-1 * conf.Mecha.RotationForceMinRange)
		} else if rotForce > 0 {
			rotForce = models.Degree(conf.Mecha.RotationForceMinRange)
		}
	} else if math.Abs(float64(rotForce)) > float64(conf.Mecha.RotationForceMaxRange) {
		if rotForce < 0 {
			rotForce = models.Degree(-1 * conf.Mecha.RotationForceMaxRange)
		} else if rotForce > 0 {
			rotForce = models.Degree(conf.Mecha.RotationForceMaxRange)
		}
	}

	force.AddRot(rotForce)
}

func (i *LookatCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	_, target := state.GetTransformByKey(i.Target)

	TockLookat(target, *i.conf, force, state)

	if math.Abs(float64(target.RobROT)) < float64(i.conf.CommandParameter.LookatToleranceDeg) {
		i.shouldClear = true
	}
}

func (i LookatCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
