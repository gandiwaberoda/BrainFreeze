package fulfillments

import (
	"fmt"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FulfillmentInterface interface {
	AsString() string
	Tick()
	ShouldClear() bool
}

var fulfillers []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface) = []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface){
	ParseHoldFulfillment,
	// ParseDurationFulfillment,
	// ParseGotballFulfillment,
	// ParseLostballFulfillment,
	// ParsePositionFulfillment,
	// ParseGlancedFulfillment,
	// ParseDistanceFulfillment,
}

func WhichFulfillment(fullbfvid string, conf *configuration.FreezeConfig, state *state.StateAccess) FulfillmentInterface {
	parsed, err := bfvid.ParseCommandSPOK(fullbfvid)
	if err != nil {
		fmt.Println("failed to parse command:", err)
		return nil
	}

	if strings.EqualFold(parsed.Fulfilment, "") {
		return nil
	}

	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(*parsed, conf, state)
		if thisIs {
			return fulfiller
		}
	}

	return nil
}
