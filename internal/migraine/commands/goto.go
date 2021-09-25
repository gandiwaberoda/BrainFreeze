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

type GotoCommand struct {
	TargetX     float64
	TargetY     float64
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseGotoCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 4 {
	// 	return false, nil
	// }

	// if !strings.EqualFold(cmd[:4], "GOTO") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "GOTO") {
		return false, nil, nil
	}

	// re, _ := regexp.Compile(`([0-9-]+),([0-9-]+)`)
	// foundParam := re.FindStringSubmatch(cmd)
	// if len(foundParam) < 3 {
	// 	fmt.Println("Nilai targetX dan targetY diperlukan")
	// 	return false, nil
	// }

	fmt.Println("GOTO PARAM::", cmd.Parameter)
	for _, v := range cmd.Parameter {
		fmt.Println("zzzz: " + v)
	}

	if len(cmd.Parameter) != 2 {
		return true, nil, errors.New("goto command require exactly 2 parameter")
	}

	tX, err := strconv.ParseFloat(cmd.Parameter[0], 64)
	if err != nil {
		return true, nil, errors.New("failed parse target X")
	}

	tY, err := strconv.ParseFloat(cmd.Parameter[1], 64)
	if err != nil {
		return true, nil, errors.New("failed parse target Y")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultPositionFulfillment(tX, tY, conf, curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := GotoCommand{
		TargetX:     tX,
		TargetY:     tY,
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i GotoCommand) GetName() string {
	return "GOTO (" + fmt.Sprint(i.TargetX) + ", " + fmt.Sprint(i.TargetY) + "):"
}

func TockGoto(targetX float64, targetY float64, conf *configuration.FreezeConfig, force *models.Force, state *state.StateAccess) {
	my := state.GetState().MyTransform

	yError := targetY - float64(my.WorldYcm)
	xError := targetX - float64(my.WorldXcm)

	sud4YRad := (my.WorldROT).AsRadian()
	yErrorRob := xError*math.Sin(float64(sud4YRad)) + yError*math.Cos(float64(sud4YRad))

	sud4XRad := (my.WorldROT * -1).AsRadian()
	xErrorRob := xError*math.Cos(float64(sud4XRad)) + yError*math.Sin(float64(sud4XRad))

	// fmt.Println("robX: " + fmt.Sprint(xErrorRob) + ";; robY: " + fmt.Sprint(yErrorRob))

	if conf.CommandParameter.AllowXYTogether {
		if math.Abs(yErrorRob) > float64(conf.CommandParameter.PositionToleranceCm) {
			force.AddY(yErrorRob)
		}

		if math.Abs(xErrorRob) > float64(conf.CommandParameter.PositionToleranceCm) {
			force.AddX(xErrorRob)
		}
	} else {
		if math.Abs(yErrorRob) > float64(conf.CommandParameter.PositionToleranceCm) {
			force.AddY(yErrorRob)
			return
		}

		if math.Abs(xErrorRob) > float64(conf.CommandParameter.PositionToleranceCm) {
			force.AddX(xErrorRob)
		}
	}
}

func (i *GotoCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	TockGoto(i.TargetX, i.TargetY, i.conf, force, state)
	force.ClampMinXY(*i.conf)

	conf := i.conf

	targetX := i.TargetX
	targetY := i.TargetY

	my := state.GetState().MyTransform
	yError := targetY - float64(my.WorldYcm)
	xError := targetX - float64(my.WorldXcm)
	// FIXME: Ini juga
	if math.Abs(yError) < float64(conf.CommandParameter.PositionToleranceCm) && math.Abs(float64(xError)) < float64(conf.CommandParameter.PositionToleranceCm) {
		i.shouldClear = true
	}
}

// func (i GotoCommand) ShouldClear() bool {
// 	return i.shouldClear
// }

func (i GotoCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
