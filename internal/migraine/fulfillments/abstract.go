package fulfillments

import (
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FulfillmentInterface interface {
	GetName() string
	Tick(*state.StateAccess) bool // return value adalah isFulfilled
}

var fulfillers []func(string, *configuration.FreezeConfig) (bool, FulfillmentInterface) = []func(string, *configuration.FreezeConfig) (bool, FulfillmentInterface){
	ParseDurationFulfillment,
}

func WhichFulfillment(intercom models.Intercom, conf *configuration.FreezeConfig) FulfillmentInterface {
	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(intercom.Content, conf)
		if thisIs {
			return fulfiller
		}
	}

	return nil
}
