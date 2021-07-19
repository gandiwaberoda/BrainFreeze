package fulfillments

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type DeltaposFuilfillment struct {
	StartTransform models.Transform
	TargetDeltaCm  models.Centimeter
	shouldClear    bool
	state          *state.StateAccess
}

// func DefaultDurationFulfillment() FulfillmentInterface {
// 	return &DurationFuilfillment{
// 		StartTime: time.Now(),
// 		Milis:     1000, // 1s
// 	}
// }

func ParseDeltaposFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "DELTAPOS") {
		return false, nil, nil
	}

	if len(fullcmd.FulfilmentParameter) != 1 {
		return true, nil, errors.New("deltapos require exactly 1 parameter which is target world position delta")
	}

	deltaworldpos, err := strconv.Atoi(fullcmd.FulfilmentParameter[0])
	if err != nil {
		return true, nil, errors.New("failed to parse parameter of deltapos fulfilment")
	}

	return true, &DeltaposFuilfillment{
		StartTransform: curstate.GetState().MyTransform,
		TargetDeltaCm:  models.Centimeter(deltaworldpos),
		state:          curstate,
	}, nil
}

func (f DeltaposFuilfillment) AsString() string {
	return fmt.Sprint("DELTAPOS (", f.TargetDeltaCm, ")")
}

func (f *DeltaposFuilfillment) Tick() {
	my := f.state.GetState().MyTransform
	start := f.StartTransform

	deltaX := my.WorldXcm - start.WorldXcm
	deltaY := my.WorldYcm - start.WorldYcm

	delta := models.Centimeter(models.EucDistance(float64(deltaX), float64(deltaY)))

	if delta >= f.TargetDeltaCm {
		f.shouldClear = true
	}
}

func (f DeltaposFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
