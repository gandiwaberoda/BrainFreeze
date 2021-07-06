package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type LostballFuilfillment struct {
	state *state.StateAccess
}

func DefaultLostballFulfillment(state *state.StateAccess) FulfillmentInterface {
	return &LostballFuilfillment{state: state}
}

func ParseLostballFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface) {
	if len(fil) < 8 || !strings.EqualFold(fil[:8], "LOSTBALL") {
		return false, nil
	}

	return true, &LostballFuilfillment{state: state}
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
