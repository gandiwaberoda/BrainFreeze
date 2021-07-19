package fulfillments

import (
	"errors"
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

func WhichFulfillment(fullbfvid string, conf *configuration.FreezeConfig, state *state.StateAccess) (FulfillmentInterface, error) {
	parsed, err := bfvid.ParseCommandSPOK(fullbfvid)
	if err != nil {
		return nil, errors.New(fmt.Sprint("failed to parse command:", err))
	}

	if strings.EqualFold(parsed.Fulfilment, "") {
		return nil, errors.New(fmt.Sprint("Fulfilment can't be of length 0"))
	}

	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(*parsed, conf, state)
		if thisIs {
			return fulfiller, nil
		}
	}

	return nil, errors.New(fmt.Sprint("fulfilment not found:", fullbfvid))
}
