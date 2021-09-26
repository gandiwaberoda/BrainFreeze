package fulfillments

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type WRotationFuilfillment struct {
	TargetRot   models.Degree
	shouldClear bool
	curstate    *state.StateAccess
	conf        *configuration.FreezeConfig

	counter int // Untuk debounce
}

func DefaultWRotationFulfillment(rotTarget models.Degree, conf *configuration.FreezeConfig, curstate *state.StateAccess) FulfillmentInterface {
	return &WRotationFuilfillment{
		TargetRot: rotTarget,
		curstate:  curstate,
		conf:      conf,
	}
}

func ParseWRotationFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "WROT") {
		return false, nil, nil
	}

	var targetRot models.Degree
	if len(fullcmd.FulfilmentParameter) == 1 {
		i, err := strconv.Atoi(fullcmd.FulfilmentParameter[0])
		if err != nil {
			return true, nil, errors.New("failed to parse parameter of wrot fulfilment: " + fullcmd.FulfilmentParameter[0])
		}
		targetRot = models.Degree(i)
	} else {
		return true, nil, errors.New("world rotation fulfilment require 1 parameter")
	}

	return true, &WRotationFuilfillment{
		TargetRot: targetRot,
		curstate:  curstate,
		conf:      conf,
	}, nil
}

func (f WRotationFuilfillment) AsString() string {
	return "WROT(" + strconv.Itoa(int(f.TargetRot)) + ")"
}

func (f *WRotationFuilfillment) Tick() {
	if math.Abs(float64(f.TargetRot-f.curstate.GetState().MyTransform.WorldROT)) < float64(f.conf.CommandParameter.LookatToleranceDeg) {
		if f.counter > 50 {
			f.shouldClear = true
		} else {
			f.counter++
		}
	} else {
		f.shouldClear = false
		f.counter = 0
	}
}

func (f WRotationFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
