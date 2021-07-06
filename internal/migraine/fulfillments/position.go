package fulfillments

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

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

func ParsePositionFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface) {
	if len(fil) < 3 || !strings.EqualFold(fil[:3], "POS") {
		return false, nil
	}

	re, _ := regexp.Compile(`([0-9-]+),([0-9-]+)`)
	foundParam := re.FindStringSubmatch(fil)
	if len(foundParam) < 3 {
		fmt.Println("Nilai targetX dan targetY diperlukan")
		return false, nil
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

	return true, &PositionFuilfillment{state: state, conf: conf, targetX: tX, targetY: tY}
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
