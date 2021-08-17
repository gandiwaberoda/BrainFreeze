package fulfillments

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GlancedFuilfillment struct {
	shouldClear bool
	state       *state.StateAccess
	targetKey   string
	conf        *configuration.FreezeConfig

	ms_to_clear         int
	lastFrameGlanced    bool `default:"false"`
	lastTimeGlanceStart time.Time
}

func DefaultGlancedFulfillment(target string, msToClear int, state *state.StateAccess, conf *configuration.FreezeConfig) FulfillmentInterface {
	return &GlancedFuilfillment{
		targetKey:   strings.ToUpper(target),
		state:       state,
		conf:        conf,
		ms_to_clear: msToClear,
	}
}

func ParseGlancedFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "GLANCED") {
		return false, nil, nil
	}

	ms := 0
	if len(fullcmd.FulfilmentParameter) == 1 {
		if found, _ := curstate.GetTransformByKey(fullcmd.FulfilmentParameter[0]); !found {
			return true, nil, errors.New("glance fulfilment target key is not recognizeable")
		}
	} else if len(fullcmd.FulfilmentParameter) == 2 {
		if found, _ := curstate.GetTransformByKey(fullcmd.FulfilmentParameter[0]); !found {
			return true, nil, errors.New("glance fulfilment target key is not recognizeable")
		}
		ms_, err := strconv.Atoi(fullcmd.FulfilmentParameter[1])
		if err != nil {
			return true, nil, errors.New("glance fulfilment ms to clear is not valid int")
		}
		ms = ms_
	} else {
		return true, nil, errors.New("glance fulfilment require 1 or 2 parameter")
	}

	return true, &GlancedFuilfillment{
		targetKey:   strings.ToUpper(fullcmd.FulfilmentParameter[0]),
		conf:        conf,
		ms_to_clear: ms,
		state:       curstate,
	}, nil
}

func (f GlancedFuilfillment) AsString() string {
	return "GLANCED(" + f.targetKey + "," + strconv.Itoa(f.ms_to_clear) + ")"
}

func (f *GlancedFuilfillment) Tick() {
	_, obj := f.state.GetTransformByKey(f.targetKey)
	if math.Abs(float64(obj.RobROT)) < float64(f.conf.CommandParameter.LookatToleranceDeg) {
		if !f.lastFrameGlanced {
			f.lastFrameGlanced = true
			f.lastTimeGlanceStart = time.Now()
			fmt.Println("LOOKEDAT: ", time.Since(f.lastTimeGlanceStart).Milliseconds())
			return
		}

		if f.lastFrameGlanced && time.Since(f.lastTimeGlanceStart).Milliseconds() >= int64(f.ms_to_clear) {
			f.shouldClear = true
		}
	} else {
		fmt.Println("NOT LOOKEDAT: ", time.Since(f.lastTimeGlanceStart).Milliseconds())
		f.shouldClear = false
		f.lastFrameGlanced = false
	}
}

func (f GlancedFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
