package fulfillments

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type DistanceFuilfillment struct {
	shouldClear bool
	state       *state.StateAccess
	targetKey   string
	dist        int
}

func DefaultDistanceFulfillment(target string, dist int, state *state.StateAccess, conf *configuration.FreezeConfig) FulfillmentInterface {
	return &DistanceFuilfillment{
		targetKey: strings.ToUpper(target),
		state:     state,
		dist:      dist,
	}
}

func ParseDistanceFulfillment(intercom models.Intercom, fil string, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface) {
	if len(fil) < 4 || !strings.EqualFold(fil[:4], "DIST") {
		return false, nil
	}

	re, _ := regexp.Compile(`\((.+),?([0-9]+)?\)`)
	foundParam := re.FindStringSubmatch(fil)
	if len(foundParam) < 2 {
		fmt.Println("Format fulfilment DIST(target, opsional rob dist)")
		return false, nil
	}

	fmt.Println(foundParam)
	for _, v := range foundParam {
		fmt.Println("x: " + v)
	}

	arg := strings.Split(foundParam[1], ",")

	defaultApproachDist := conf.CommandParameter.ApproachDistanceCm

	var target string
	var dist int

	if len(arg) == 1 {
		target = strings.ToUpper(arg[0])
		dist = defaultApproachDist
	} else if len(arg) == 2 {
		target = strings.ToUpper(arg[0])

		if _parsed, err := strconv.Atoi(arg[1]); err == nil {
			fmt.Println("C")
			if _parsed == 0 {
				dist = defaultApproachDist
			} else {
				dist = _parsed
			}
		} else {
			dist = defaultApproachDist
		}
	}

	return true, &DistanceFuilfillment{
		targetKey: target,
		dist:      dist,
		state:     state,
	}
}

func (f DistanceFuilfillment) AsString() string {
	return "DIST(" + f.targetKey + "," + fmt.Sprint(f.dist) + ")"
}

func (f *DistanceFuilfillment) Tick() {
	_, obj := f.state.GetTransformByKey(f.targetKey)
	if math.Abs(float64(obj.RobRcm)) < float64(f.dist) {
		f.shouldClear = true
	}
}

func (f DistanceFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
