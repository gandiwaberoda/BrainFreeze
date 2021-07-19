package fulfillments

import (
	"errors"
	"math"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GlancedFuilfillment struct {
	shouldClear bool
	state       *state.StateAccess
	targetKey   string
	conf        *configuration.FreezeConfig
}

func DefaultGlancedFulfillment(target string, state *state.StateAccess, conf *configuration.FreezeConfig) FulfillmentInterface {
	return &GlancedFuilfillment{
		targetKey: strings.ToUpper(target),
		state:     state,
		conf:      conf,
	}
}

func ParseGlancedFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "GLANCED") {
		return false, nil, nil
	}

	if len(fullcmd.FulfilmentParameter) == 1 {
		if found, _ := state.GetTransformByKey(fullcmd.FulfilmentParameter[0]); !found {
			return true, nil, errors.New("glance fulfilment target key is not recognizeable")
		}
	} else {
		return true, nil, errors.New("glance fulfilment require exactly 1 parameter")
	}

	return true, &GlancedFuilfillment{
		targetKey: fullcmd.FulfilmentParameter[0],
	}, nil
}

func (f GlancedFuilfillment) AsString() string {
	return "GLANCED(" + f.targetKey + ")"
}

func (f *GlancedFuilfillment) Tick() {
	_, obj := f.state.GetTransformByKey(f.targetKey)
	if math.Abs(float64(obj.RobROT)) < float64(f.conf.CommandParameter.LookatToleranceDeg) {
		f.shouldClear = true
	} else {
		f.shouldClear = false
	}
}

func (f GlancedFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
