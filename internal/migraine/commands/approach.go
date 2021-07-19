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

type ApproachCommand struct {
	Target      string
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseApproachCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 8 || !strings.EqualFold(cmd[:8], "APPROACH") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "APPROACH") {
		return false, nil, nil
	}

	// re, _ := regexp.Compile(`\((.+)\)`)
	// foundParam := re.FindStringSubmatch(cmd)
	// if len(foundParam) < 1 {
	// 	fmt.Println("Target diperlukan")
	// 	return false, nil
	// }
	// fmt.Println(foundParam)
	// target := strings.ToUpper(foundParam[1])
	if len(cmd.Parameter) != 1 {
		return true, nil, errors.New("approach only accept one parameter")
	}
	target := cmd.Parameter[0]

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		fmt.Println("Target key not acceptable")
		return true, nil, errors.New("key is not acceptable")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultDistanceFulfillment(target, conf.CommandParameter.ApproachDistanceCm, curstate, conf)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := ApproachCommand{
		Target:      target,
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i ApproachCommand) GetName() string {
	return "APPROACH (" + string(i.Target) + ")"
}

func TockApproach(target models.Transform, conf configuration.FreezeConfig, force *models.Force, state *state.StateAccess) {
	xF := target.RobXcm
	yF := target.RobYcm

	if conf.CommandParameter.AllowXYTogether {
		force.AddX(float64(xF))
		force.AddY(float64(yF))
	} else {
		force.AddY(float64(yF))
	}
}

func (i *ApproachCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	_, target := state.GetTransformByKey(i.Target)

	TockApproach(target, *i.conf, force, state)
}

func (i ApproachCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
