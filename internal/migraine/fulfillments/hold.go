package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type HoldFuilfillment struct {
}

func DefaultHoldFulfillment() *HoldFuilfillment {
	return &HoldFuilfillment{}
}

func ParseHoldFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig) (bool, FulfillmentInterface) {
	if !strings.EqualFold(fil[:4], "HOLD") {
		return false, nil
	}

	return true, &HoldFuilfillment{}
}

func (f HoldFuilfillment) AsString() string {
	return "HOLD"
}

func (f *HoldFuilfillment) Tick(state *state.StateAccess) bool {
	return false
}

func (f HoldFuilfillment) ShouldClear() bool {
	return false
}
