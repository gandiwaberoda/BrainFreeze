package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type LostballFuilfillment struct {
	state *state.StateAccess
}

func DefaultLostballFulfillment(state *state.StateAccess) FulfillmentInterface {
	return &LostballFuilfillment{state: state}
}

func ParseLostballFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "LOSTBALL") {
		return false, nil, nil
	}

	return true, &LostballFuilfillment{state: state}, nil
}

func (f LostballFuilfillment) AsString() string {
	return "LOSTBALL"
}

func (f *LostballFuilfillment) Tick() {
}

func (f LostballFuilfillment) ShouldClear() bool {
	// TODO: Buat supaya flag berubah setelah isDribbling false sekian milisecond (Debounce)
	return !f.state.GetState().GutToBrain.IsDribbling
}
