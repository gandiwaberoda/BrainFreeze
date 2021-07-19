package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type GotballFuilfillment struct {
	state *state.StateAccess
}

func DefaultGotballFulfillment(state *state.StateAccess) FulfillmentInterface {
	return &GotballFuilfillment{state: state}
}

func ParseGotballFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "GOTBALL") {
		return false, nil, nil
	}

	return true, &GotballFuilfillment{state: state}, nil
}

func (f GotballFuilfillment) AsString() string {
	return "GOTBALL"
}

func (f *GotballFuilfillment) Tick() {
}

func (f GotballFuilfillment) ShouldClear() bool {
	return f.state.GetState().GutToBrain.IsDribbling
}
