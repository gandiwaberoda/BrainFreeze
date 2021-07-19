package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type HoldFuilfillment struct {
}

func DefaultHoldFulfillment() FulfillmentInterface {
	return &HoldFuilfillment{}
}

func ParseHoldFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "HOLD") {
		return false, nil, nil
	}

	return true, &HoldFuilfillment{}, nil
}

func (f HoldFuilfillment) AsString() string {
	return "HOLD"
}

func (f *HoldFuilfillment) Tick() {
}

func (f HoldFuilfillment) ShouldClear() bool {
	return false
}
