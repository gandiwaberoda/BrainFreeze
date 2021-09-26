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

var fulfillers []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface, error) = []func(bfvid.CommandSPOK, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface, error){
	ParseHoldFulfillment,
	ParseDurationFulfillment,
	ParseGotballFulfillment,
	ParseLostballFulfillment,
	ParsePositionFulfillment,
	ParseGlancedFulfillment,
	ParseDistanceFulfillment,
	ParseDeltaposFulfillment,
	ParseFrontFulfillment,
	ParseWRotationFulfillment,
}

func WhichFulfillment(fullbfvid string, conf *configuration.FreezeConfig, state *state.StateAccess) (FulfillmentInterface, error) {
	parsed, err := bfvid.ParseCommandSPOK(fullbfvid)
	if err != nil {
		return nil, errors.New(fmt.Sprint("failed to parse command:", err))
	}

	if strings.EqualFold(parsed.Fulfilment, "") {
		return nil, errors.New("fulfilment can't be of length 0")
	}

	for _, isThis := range fulfillers {
		thisIs, fulfiller, err := isThis(*parsed, conf, state)
		if thisIs {
			if err != nil {
				return nil, errors.New(fmt.Sprint("failed to parse fulfilment:", err))
			}
			return fulfiller, nil
		}
	}

	return nil, errors.New(fmt.Sprint("fulfilment not found:", parsed))
}
