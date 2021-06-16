package fulfillments

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type DurationFuilfillment struct {
	StartTime time.Time
	Milis     models.Miliseconds
}

func DefaultDurationFulfillment() DurationFuilfillment {
	return DurationFuilfillment{
		StartTime: time.Now(),
		Milis:     1000, // 1s
	}
}

func ParseDurationFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig) (bool, FulfillmentInterface) {
	if !strings.EqualFold(fil[:3], "DUR") {
		return false, nil
	}

	re, _ := regexp.Compile("([0-9]+)")
	foundParam := re.FindString(fil)

	var milis models.Miliseconds
	if foundParam != "" {
		i, err := strconv.Atoi(foundParam)
		if err != nil {
			fmt.Println("Failed parsing ATOI:", err)
			milis = models.Miliseconds(conf.Fulfillment.DefaultDurationMs)
		} else {
			milis = models.Miliseconds(i)
		}
	} else {
		milis = models.Miliseconds(conf.Fulfillment.DefaultDurationMs)
	}

	return true, DurationFuilfillment{
		StartTime: time.Now(),
		Milis:     models.Miliseconds(milis),
	}
}

func (f DurationFuilfillment) AsString() string {
	return "DUR(" + strconv.Itoa(int(f.Milis)) + ")"
}

func (f DurationFuilfillment) Tick(state *state.StateAccess) bool {
	elapsed := time.Since(f.StartTime)
	return elapsed.Milliseconds() > int64(f.Milis)
}
