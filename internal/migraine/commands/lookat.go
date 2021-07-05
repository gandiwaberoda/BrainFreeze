package commands

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

// TODO: BUATKAN LOOKEDAT Fulfillment
type LookatCommand struct {
	Target      string
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseLookatCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 6 {
		return false, nil
	}

	if !strings.EqualFold(cmd[:6], "LOOKAT") {
		return false, nil
	}

	re, _ := regexp.Compile(`\(([A-Za-z0-9]+)\)`)
	foundParam := re.FindString(cmd)
	foundParam = strings.ReplaceAll(foundParam, "(", "")
	foundParam = strings.ReplaceAll(foundParam, ")", "")

	fmt.Println("zzz", foundParam)

	target := "BALL"
	if foundParam != "" {
		fmt.Println(foundParam)
		target = foundParam
	}

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultComplexFulfillment()
	}
	parsed := LookatCommand{
		Target:      target,
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed
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
	_, target := state.GetTransformByKey(i.Target)

	TockLookat(target, *i.conf, force, state)

	// FIXME: Ini perlu ganti pake fulfillmentnya tersendiri
	if math.Abs(float64(target.RobROT)) < float64(i.conf.CommandParameter.LookatToleranceDeg) {
		i.shouldClear = true
	}
	i.fulfillment.Tick()
}

func (i LookatCommand) ShouldClear() bool {
	return i.shouldClear
}

func (i LookatCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
