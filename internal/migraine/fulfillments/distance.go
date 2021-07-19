package fulfillments

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
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

func ParseDistanceFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "DIST") {
		return false, nil, nil
	}

	// re, _ := regexp.Compile(`\((.+),?([0-9]+)?\)`)
	// foundParam := re.FindStringSubmatch(fil)
	// if len(foundParam) < 2 {
	// 	fmt.Println("Format fulfilment DIST(target, opsional rob dist)")
	// 	return false, nil
	// }

	// fmt.Println(foundParam)
	// for _, v := range foundParam {
	// 	fmt.Println("x: " + v)
	// }

	// arg := strings.Split(foundParam[1], ",")
	arg := fullcmd.FulfilmentParameter

	defaultApproachDist := conf.CommandParameter.ApproachDistanceCm

	var target string
	var dist int

	if len(arg) == 1 {
		target = strings.ToUpper(arg[0])
		dist = defaultApproachDist
	} else if len(arg) == 2 {
		target = strings.ToUpper(arg[0])

		if _parsed, err := strconv.Atoi(arg[1]); err == nil {
			if _parsed == 0 {
				dist = defaultApproachDist
			} else {
				dist = _parsed
			}
		} else {
			return true, nil, errors.New(fmt.Sprint("distance fulfilment second parameter of (", arg[1], ") is failed to int"))
		}
	} else {
		return true, nil, errors.New("distance fulfilment require either 1 or 2 parameter")
	}

	if !state.GetTransformKeyAcceptable(target) {
		return true, nil, errors.New(fmt.Sprint("Target of (", target, ") is not recognizeable"))
	}

	return true, &DistanceFuilfillment{
		targetKey: target,
		dist:      dist,
		state:     curstate,
	}, nil
}

func (f DistanceFuilfillment) AsString() string {
	return "DIST(" + f.targetKey + "," + fmt.Sprint(f.dist) + ")"
}

func (f *DistanceFuilfillment) Tick() {
	_, obj := f.state.GetTransformByKey(f.targetKey)

	// FIXME: Pake value dari expired, kalau expired berarti belum terpenuhi

	if math.Abs(float64(obj.RobRcm)) < float64(f.dist) {
		f.shouldClear = true
	}
}

func (f DistanceFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
