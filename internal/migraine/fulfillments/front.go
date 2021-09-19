/*
Lihat warna apa yang ada di garis didepan, abaikan warna oranye bola
(Basically cuma lihat MAGENTA, BIRU, MERAH), skip ORANGE <- cuma MAGENTA, BIRU dan MERAH yang bisa mentrigger ultrasonicnya (ANGGAP)
*/

package fulfillments

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type FrontFuilfillment struct {
	state        *state.StateAccess
	conf         *configuration.FreezeConfig
	targetColor  string
	targetDistCm int
	shouldClear  bool
}

func DefaultFrontFulfillment(targetColor string, targetDistCm int, conf *configuration.FreezeConfig, state *state.StateAccess) FulfillmentInterface {
	return &FrontFuilfillment{state: state, conf: conf, targetColor: targetColor, targetDistCm: targetDistCm}
}

func ParseFrontFulfillment(fullcmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, state *state.StateAccess) (bool, FulfillmentInterface, error) {
	if !strings.EqualFold(fullcmd.Fulfilment, "FRONT") {
		return false, nil, nil
	}

	if len(fullcmd.FulfilmentParameter) != 2 {
		return true, nil, errors.New("front fulfilment require exactly 2 parameter, color and distcm")
	}

	colTarget := fullcmd.FulfilmentParameter[0]
	acc := false
	for _, v := range []string{"CYAN", "MAGENTA", "DUMMY"} {
		if strings.EqualFold(v, colTarget) {
			acc = true
		}
	}
	if !acc {
		return true, nil, errors.New("target color of front fulfilment isn't valid")
	}

	distTarget, errDist := strconv.ParseInt(fullcmd.FulfilmentParameter[1], 10, 64)
	if errDist != nil {
		return true, nil, errors.New("failed to parse targetdistcm of post fulfilment")
	}

	return true, &FrontFuilfillment{state: state, conf: conf, targetColor: colTarget, targetDistCm: int(distTarget)}, nil
}

func (f FrontFuilfillment) AsString() string {
	return "FRONT(" + fmt.Sprint(f.targetColor) + "," + fmt.Sprint(f.targetDistCm) + ")"
}

func (f *FrontFuilfillment) Tick() {
	frontDist := f.state.GetState().Araya.Dist0

	acc := false
	for _, v := range f.state.GetState().Straight {
		if strings.EqualFold("BALL", v.DetectedColorName) {
			continue // Abaikan bola
		}

		if strings.EqualFold(v.DetectedColorName, f.targetColor) {
			acc = true
			break
		}

		if !strings.EqualFold(v.DetectedColorName, f.targetColor) {
			acc = false
			break
		}
	}

	if acc {
		if frontDist < models.Centimeter(f.targetDistCm) {
			f.shouldClear = true
			return
		}
		fmt.Println("A", frontDist)
	}

	fmt.Println("B")
	f.shouldClear = false
}

func (f FrontFuilfillment) ShouldClear() bool {
	return f.shouldClear
}
