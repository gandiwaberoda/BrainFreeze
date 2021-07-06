package fulfillments

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
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

func ParseGlancedFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface) {
	if !strings.EqualFold(fil[:7], "GLANCED") {
		return false, nil
	}

	re, _ := regexp.Compile("(.)+")
	foundParam := re.FindString(fil)

	if foundParam != "" {
		if found, _ := state.GetTransformByKey(foundParam); !found {
			fmt.Println("Target key is not recognizeable")
			return false, nil
		}
	} else {
		fmt.Println("No glance argument found")
		return false, nil
	}

	return true, &GlancedFuilfillment{
		targetKey: foundParam,
	}
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
