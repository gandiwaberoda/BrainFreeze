package fulfillments

import (
	"strconv"
)

type ComplexFuilfillment struct {
	fulfilled bool
}

func DefaultComplexFulfillment() FulfillmentInterface {
	return &ComplexFuilfillment{
		fulfilled: false,
	}
}

func (f *ComplexFuilfillment) Fulfilled() {
	f.fulfilled = true
}

// func (f ComplexFuilfillment) IsFulfilled() bool {
// 	return f.fulfilled
// }

func (f ComplexFuilfillment) AsString() string {
	return "COMPLEX(" + strconv.FormatBool(f.ShouldClear()) + ")"
}

func (f *ComplexFuilfillment) Tick() {
}

func (f ComplexFuilfillment) ShouldClear() bool {
	return f.fulfilled
}
