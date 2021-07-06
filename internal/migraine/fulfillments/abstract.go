package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FulfillmentInterface interface {
	AsString() string
	Tick() // return value adalah isFulfilled
	ShouldClear() bool
}

var fulfillers []func(models.Intercom, string, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface) = []func(models.Intercom, string, *configuration.FreezeConfig, *state.StateAccess) (bool, FulfillmentInterface){
	ParseDurationFulfillment,
	ParseGotballFulfillment,
	ParseLostballFulfillment,
	ParsePositionFulfillment,
	ParseHoldFulfillment,
}

func WhichFulfillment(intercom models.Intercom, conf *configuration.FreezeConfig, state *state.StateAccess) FulfillmentInterface {
	splitted := strings.Split(intercom.Content, "/")

	if len(splitted) < 2 {
		return nil
	}

	if strings.EqualFold(splitted[1], "") {
		return nil
	}

	filMsg := strings.ToUpper(strings.TrimSpace(splitted[1]))

	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(intercom, filMsg, conf, state)
		if thisIs {
			return fulfiller
		}
	}

	return nil
}
