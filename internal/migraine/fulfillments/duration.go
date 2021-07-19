package fulfillments

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type DurationFuilfillment struct {
	StartTime   time.Time
	Milis       models.Miliseconds
	shouldClear bool
	elapsed     time.Duration
}

func DefaultDurationFulfillment() FulfillmentInterface {
	return &DurationFuilfillment{
		StartTime: time.Now(),
		Milis:     1000, // 1s
	}
}

func ParseDurationFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "DUR") {
		return false, nil, nil
	}

	fmt.Println(fullcmd)

	var milis models.Miliseconds
	if len(fullcmd.FulfilmentParameter) == 1 {
		i, err := strconv.Atoi(fullcmd.Fulfilment)
		if err != nil {
			milis = models.Miliseconds(i)
		} else {
			return true, nil, errors.New("failed to parse parameter of duration fulfilment")
		}
	} else if len(fullcmd.FulfilmentParameter) == 0 {
		milis = models.Miliseconds(conf.Fulfillment.DefaultDurationMs)
	} else {
		return true, nil, errors.New("duration fulfilment require either 1 or none parameter")
	}

	return true, &DurationFuilfillment{
		StartTime: time.Now(),
		Milis:     models.Miliseconds(milis),
	}, nil
}

func (f DurationFuilfillment) AsString() string {
	return "DUR(" + strconv.Itoa(int(f.Milis)) + ")"
}

func (f *DurationFuilfillment) Tick() {
	elapsed := time.Since(f.StartTime)
	fulfilled := elapsed.Milliseconds() > int64(f.Milis)
	f.shouldClear = fulfilled
	f.elapsed = elapsed
}

func (f DurationFuilfillment) ShouldClear() bool {
	return f.shouldClear
}

func (f DurationFuilfillment) GetElapsed() time.Duration {
	return f.elapsed
}
