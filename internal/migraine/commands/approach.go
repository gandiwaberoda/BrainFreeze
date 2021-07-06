package commands

import (
	"fmt"
	"regexp"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
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
func ParseApproachCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 8 || !strings.EqualFold(cmd[:8], "APPROACH") {
		return false, nil
	}

	re, _ := regexp.Compile(`\((.+)\)`)
	foundParam := re.FindStringSubmatch(cmd)
	if len(foundParam) < 1 {
		fmt.Println("Target diperlukan")
		return false, nil
	}
	fmt.Println(foundParam)
	target := strings.ToUpper(foundParam[1])

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		fmt.Println("Target key not acceptable")
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultDistanceFulfillment(target, conf.CommandParameter.ApproachDistanceCm, curstate, conf)
	}
	parsed := ApproachCommand{
		Target:      target,
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed
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
	_, target := state.GetTransformByKey(i.Target)

	TockApproach(target, *i.conf, force, state)

	i.fulfillment.Tick()
}

func (i ApproachCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
