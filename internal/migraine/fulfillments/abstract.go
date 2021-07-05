package fulfillments

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FulfillmentInterface interface {
	AsString() string
	Tick(*state.StateAccess) bool // return value adalah isFulfilled
	ShouldClear() bool
}

var fulfillers []func(models.Intercom, string, *configuration.FreezeConfig) (bool, FulfillmentInterface) = []func(models.Intercom, string, *configuration.FreezeConfig) (bool, FulfillmentInterface){
	ParseDurationFulfillment,
}

func WhichFulfillment(intercom models.Intercom, conf *configuration.FreezeConfig) FulfillmentInterface {
	splitted := strings.Split(intercom.Content, "/")

	if len(splitted) < 2 {
		return nil
	}

	if strings.EqualFold(splitted[1], "") {
		return nil
	}

	filMsg := strings.ToUpper(strings.TrimSpace(splitted[1]))

	for _, isThis := range fulfillers {
		thisIs, fulfiller := isThis(intercom, filMsg, conf)
		if thisIs {
			return fulfiller
		}
	}

	return nil
}
