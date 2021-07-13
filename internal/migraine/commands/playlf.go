package commands

import (
	"math"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PlaylfCommand struct {
	fulfillment      fulfillments.FulfillmentInterface
	conf             *configuration.FreezeConfig
	lastRotationTime time.Time
}

func ParsePlaylfCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 6 || strings.ToUpper(cmd[:6]) != "PLAYLF" {
		return false, &PlaylfCommand{}
	}

	parsedFulfillment := fulfillments.WhichFulfillment(intercom, conf, curstate)
	if parsedFulfillment == nil {
		parsedFulfillment = fulfillments.DefaultHoldFulfillment()
	}

	parsed := PlaylfCommand{
		fulfillment:      parsedFulfillment,
		conf:             conf,
		lastRotationTime: time.Now(),
	}

	return true, &parsed
}

func (i PlaylfCommand) GetName() string {
	return "PLAYLF"
}

func (i *PlaylfCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	circles := state.GetState().CircularFieldLine

	if len(circles) < 1 {
		return
	}

	// Rotasi
	rotError := circles[0]
	if models.Degree(math.Abs(rotError)) > models.Degree(i.conf.CommandParameter.LookatToleranceDeg) {
		i.lastRotationTime = time.Now()
		force.AddRot(models.Degree(rotError))

		if !i.conf.CommandParameter.AllowXYRotTogether {
			// Jika hanya boleh satu degree dalam satu waktu
			// Rotasi rotasi aja dulu
			return
		}
	}

	if time.Since(i.lastRotationTime) < time.Duration(i.conf.CommandParameter.RotToMoveDelay) && !i.conf.CommandParameter.AllowXYRotTogether {
		// Kasih delay ketika berpindah dari rotasi ke translasi
		force.Idle()
		return
	}

	robRcm := models.PxToCm(150)
	robROTDegree := models.Degree(rotError)
	robXcm := robRcm * math.Sin(float64(robROTDegree.AsRadian()))
	robYcm := robRcm * math.Cos(float64(robROTDegree.AsRadian()))

	// Bergerak
	if i.conf.CommandParameter.AllowXYTogether {
		force.AddX(robXcm)
		force.AddY(robYcm)
	} else {
		force.AddY(robYcm)
	}

}

func (i PlaylfCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}