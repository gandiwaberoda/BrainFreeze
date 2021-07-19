package commands

import (
	"math"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

// TODO: BUATKAN LOOKEDAT Fulfillment
type GetballCommand struct {
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
	// shouldClear      bool
	lastRotationTime time.Time
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseGetballCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 7 {
	// 	return false, nil
	// }

	// if !strings.EqualFold(cmd[:7], "GETBALL") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "GETBALL") {
		return false, nil, nil
	}

	parseFulfilment := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
	if parseFulfilment == nil {
		parseFulfilment = fulfillments.DefaultGotballFulfillment(curstate)
	}

	parsed := GetballCommand{
		conf:             conf,
		fulfillment:      parseFulfilment,
		lastRotationTime: time.Now(),
	}

	return true, &parsed, nil
}

func (i GetballCommand) GetName() string {
	return "GETBALL"
}

// return value adalah ShouldUpdateLastRot
func TockGetball(lastRotTime time.Time, conf configuration.FreezeConfig, force *models.Force, state *state.StateAccess) bool {
	ballState := state.GetState().BallTransform

	// Rotasi
	rotError := state.GetState().BallTransform.RobROT
	if models.Degree(math.Abs(float64(rotError))) > models.Degree(conf.CommandParameter.LookatToleranceDeg) {
		TockLookat(ballState, conf, force, state)

		if !conf.CommandParameter.AllowXYRotTogether {
			// Jika hanya boleh satu degree dalam satu waktu
			// Rotasi rotasi aja dulu
			return true
		}
	}

	if time.Since(lastRotTime) < time.Duration(conf.CommandParameter.RotToMoveDelay) && !conf.CommandParameter.AllowXYRotTogether {
		// Kasih delay ketika berpindah dari rotasi ke translasi
		force.Idle()
		return false
	}

	// Handling
	// TODO: Pake jarak CM bukan PX
	if ballState.TopRpx <= models.Centimeter(conf.CommandParameter.HandlingOnDist) {
		force.EnableHandling()
	}

	limitedY := float64(ballState.RobYcm)
	if limitedY > float64(conf.Mecha.VerticalForceRange) {
		limitedY = float64(conf.Mecha.VerticalForceRange)
	}
	if limitedY < float64(-1*conf.Mecha.VerticalForceRange) {
		limitedY = float64(-1 * conf.Mecha.VerticalForceRange)
	}

	limitedX := float64(ballState.RobXcm)
	if limitedX > float64(conf.Mecha.HorizontalForceRange) {
		limitedX = float64(conf.Mecha.HorizontalForceRange)
	}
	if limitedX < float64(-1*conf.Mecha.HorizontalForceRange) {
		limitedX = float64(-1 * conf.Mecha.HorizontalForceRange)
	}

	// Dekati bola kedepan
	if conf.CommandParameter.AllowXYTogether {
		force.AddX(limitedX)
		force.AddY(limitedY)
	} else {
		force.AddY(limitedY)
	}

	return false
}

func (i *GetballCommand) Tick(force *models.Force, state *state.StateAccess) {
	shouldUpdateLastRotTime := TockGetball(i.lastRotationTime, *i.conf, force, state)
	if shouldUpdateLastRotTime {
		i.lastRotationTime = time.Now()
	}

	i.fulfillment.Tick()

	// if state.GetState().GutToBrain.IsDribbling {
	// 	// i.shouldClear = true
	// 	i.fulfillment.Fulfilled()
	// }
}

// func (i GetballCommand) ShouldClear() bool {
// 	return i.shouldClear
// }

func (i GetballCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
