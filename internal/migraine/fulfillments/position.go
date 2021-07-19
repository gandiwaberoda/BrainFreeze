package fulfillments

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PositionFuilfillment struct {
	state       *state.StateAccess
	conf        *configuration.FreezeConfig
	targetX     int
	targetY     int
	shouldClear bool
}

func DefaultPositionFulfillment(xTarget, yTarget int, conf *configuration.FreezeConfig, state *state.StateAccess) FulfillmentInterface {
	return &PositionFuilfillment{state: state, conf: conf, targetX: xTarget, targetY: yTarget}
}

func ParsePositionFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "POS") {
		return false, nil, nil
	}

	if len(fullcmd.Parameter) != 2 {
		return true, nil, errors.New("position fulfilment require exactly 2 parameter")
	}

	tX, errX := strconv.Atoi(fullcmd.Parameter[0])
	if errX != nil {
		return true, nil, errors.New("failed to parse targetX of Position fulfilment")
	}
	tY, errY := strconv.Atoi(fullcmd.Parameter[1])
	if errY != nil {
		return true, nil, errors.New("failed to parse targetY of Position fulfilment")
	}

	return true, &PositionFuilfillment{state: state, conf: conf, targetX: tX, targetY: tY}, nil
}

func (f PositionFuilfillment) AsString() string {
	return "POS(" + fmt.Sprint(f.targetX) + "," + fmt.Sprint(f.targetY) + ")"
}

func (f *PositionFuilfillment) Tick() {
	xErr := math.Abs(float64(f.state.GetState().MyTransform.WorldXcm - models.Centimeter(f.targetX)))
	yErr := math.Abs(float64(f.state.GetState().MyTransform.WorldYcm - models.Centimeter(f.targetY)))

	if xErr > float64(f.conf.CommandParameter.PositionToleranceCm) || yErr > float64(f.conf.CommandParameter.PositionToleranceCm) {
		f.shouldClear = false
		return
	}

	f.shouldClear = true
}

func (f PositionFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
