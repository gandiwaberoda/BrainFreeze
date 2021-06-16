package fulfillments

import (
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FulfillmentInterface interface {
	GetName() string
	Tick(*state.StateAccess) bool // return value adalah isFulfilled
}
