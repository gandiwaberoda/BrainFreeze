package commands

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
	if len(cmd.Parameter) == 1 && cmd.Parameter[0] != "" {
		target = cmd.Parameter[0]
	} else if len(cmd.Parameter) > 1 {
		return true, nil, errors.New("lookat command require either NONE or 1 parameter")
	}

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		return true, nil, errors.New("lookat target key is not acceptable")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultGlancedFulfillment(target, curstate, conf)
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
	if rotForce < models.Degree(-1*conf.Mecha.RotationForceRange) {
		rotForce = models.Degree(-1 * conf.Mecha.RotationForceRange)
	}

	if rotForce > models.Degree(conf.Mecha.RotationForceRange) {
		rotForce = models.Degree(conf.Mecha.RotationForceRange)
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
