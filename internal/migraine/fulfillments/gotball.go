package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GotballFuilfillment struct {
	state *state.StateAccess
}

func DefaultGotballFulfillment(state *state.StateAccess) FulfillmentInterface {
	return &GotballFuilfillment{state: state}
}

func ParseGotballFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface) {
	if len(fil) < 7 || !strings.EqualFold(fil[:7], "GOTBALL") {
		return false, nil
	}

	return true, &GotballFuilfillment{state: state}
}

func (f GotballFuilfillment) AsString() string {
	return "GOTBALL"
}

func (f *GotballFuilfillment) Tick() {
}

func (f GotballFuilfillment) ShouldClear() bool {
	return f.state.GetState().GutToBrain.IsDribbling
}
