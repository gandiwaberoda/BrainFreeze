package parser

import (
	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

var fulfillers []func(string, *configuration.FreezeConfig) (bool, fulfillments.FulfillmentInterface) = []func(string, *configuration.FreezeConfig) (bool, fulfillments.FulfillmentInterface){
	fulfillments.ParseDurationFulfillment,
}

func WhichFulfillment(intercom models.Intercom, conf *configuration.FreezeConfig) fulfillments.FulfillmentInterface {
	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(intercom.Content, conf)
		if thisIs {
			return fulfiller
		}
	}

	return nil
}
