package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GotoCommand struct {
	TargetX     int
	TargetY     int
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
	shouldClear bool
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseGotoCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 4 {
		return false, nil
	}

	if !strings.EqualFold(cmd[:4], "GOTO") {
		return false, nil
	}

	re, _ := regexp.Compile(`([0-9-]+),([0-9-]+)`)
	foundParam := re.FindStringSubmatch(cmd)
	if len(foundParam) < 3 {
		fmt.Println("Nilai targetX dan targetY diperlukan")
		return false, nil
	}

	fmt.Println(foundParam)
	for _, v := range foundParam {
		fmt.Println("x: " + v)
	}

	tX, errX := strconv.Atoi(foundParam[1])
	if errX != nil {
		fmt.Println("Failed parse target X")
		return false, nil
	}
	tY, errY := strconv.Atoi(foundParam[2])
	if errY != nil {
		fmt.Println("Failed parse target Y")
		return false, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultPositionFulfillment(tX, tY, conf, curstate)
	}

	parsed := GotoCommand{
		TargetX:     tX,
		TargetY:     tY,
		conf:        conf,
		fulfillment: parseFulfilment,
	}

	return true, &parsed
}

func (i GotoCommand) GetName() string {
	return "GOTO (" + fmt.Sprint(i.TargetX) + ", " + fmt.Sprint(i.TargetY) + "):"
}

func TockGoto(targetX int, targetY int, conf *configuration.FreezeConfig, force *models.Force, state *state.StateAccess) {
	my := state.GetState().MyTransform

	yError := float64(targetY - int(my.WorldYcm))
	xError := float64(targetX - int(my.WorldXcm))

	sud4YRad := (my.WorldROT).AsRadian()
	yErrorRob := xError*math.Sin(float64(sud4YRad)) + yError*math.Cos(float64(sud4YRad))

	sud4XRad := (my.WorldROT * -1).AsRadian()
	xErrorRob := xError*math.Cos(float64(sud4XRad)) + yError*math.Sin(float64(sud4XRad))

	// fmt.Println("robX: " + fmt.Sprint(xErrorRob) + ";; robY: " + fmt.Sprint(yErrorRob))

	if conf.CommandParameter.AllowXYTogether {
		if int(math.Abs(yErrorRob)) > conf.CommandParameter.PositionToleranceCm {
			force.AddY(yErrorRob)
		}

		if int(math.Abs(xErrorRob)) > conf.CommandParameter.PositionToleranceCm {
			force.AddX(xErrorRob)
		}
	} else {
		if int(math.Abs(yErrorRob)) > conf.CommandParameter.PositionToleranceCm {
			force.AddY(yErrorRob)
			return
		}

		if int(math.Abs(xErrorRob)) > conf.CommandParameter.PositionToleranceCm {
			force.AddX(xErrorRob)
		}
	}
}

func (i *GotoCommand) Tick(force *models.Force, state *state.StateAccess) {
	TockGoto(i.TargetX, i.TargetY, i.conf, force, state)

	conf := i.conf

	targetX := i.TargetX
	targetY := i.TargetY

	my := state.GetState().MyTransform
	yError := float64(targetY - int(my.WorldYcm))
	xError := float64(targetX - int(my.WorldXcm))
	// FIXME: Ini juga
	if math.Abs(yError) < float64(conf.CommandParameter.PositionToleranceCm) && math.Abs(float64(xError)) < float64(conf.CommandParameter.PositionToleranceCm) {
		i.shouldClear = true
	}

	i.fulfillment.Tick()
}

// func (i GotoCommand) ShouldClear() bool {
// 	return i.shouldClear
// }

func (i GotoCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
