package fulfillments

import (
	"strconv"

	"harianugrah.com/brainfreeze/pkg/models/state"
)

type ComplexFuilfillment struct {
	fulfilled bool
}

func DefaultComplexFulfillment() ComplexFuilfillment {
	return ComplexFuilfillment{
		fulfilled: false,
	}
}

func (f *ComplexFuilfillment) Fulfilled() {
	f.fulfilled = true
}

func (f ComplexFuilfillment) IsFulfilled() bool {
	return f.fulfilled
}

func (f ComplexFuilfillment) AsString() string {
	return "COMPLEX(" + strconv.FormatBool(f.IsFulfilled()) + ")"
}

func (f ComplexFuilfillment) Tick(state *state.StateAccess) bool {
	return f.IsFulfilled()
}
