package commands

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

// Look rotation in world space
type LookWRotCommand struct {
	Target      models.Degree
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// LookWRotCommand memiliki fulfillment default wrot
func ParseLookWRotCommandCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "LOOKWROT") {
		return false, nil, nil
	}

	degTarget := models.Degree(0)
	if len(cmd.Parameter) == 1 && cmd.Parameter[0] != "" {
		_t, _err := strconv.Atoi(cmd.Parameter[0])
		if _err != nil {
			return true, nil, errors.New("LookWRotCommand can't convert " + cmd.Parameter[0] + " to int")
		}
		degTarget = models.Degree(_t)
	} else {
		return true, nil, errors.New("LookWRotCommand require exactly 1 argument")
	}
	degTarget = degTarget.AsHalfCircle()

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultWRotationFulfillment(degTarget, conf, curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := LookWRotCommand{
		Target:      degTarget,
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i LookWRotCommand) GetName() string {
	return fmt.Sprint("LookWRot:", int(i.Target))
}

func (i *LookWRotCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	rotForce := i.Target.ShiftRight() - state.GetState().MyTransform.WorldROT.ShiftRight()

	if math.Abs(float64(rotForce)) < float64(i.conf.Mecha.RotationForceMinRange) {
		if rotForce < 0 {
			rotForce = models.Degree(-1 * i.conf.Mecha.RotationForceMinRange)
		} else if rotForce > 0 {
			rotForce = models.Degree(i.conf.Mecha.RotationForceMinRange)
		}
	} else if math.Abs(float64(rotForce)) > float64(i.conf.Mecha.RotationForceMaxRange) {
		if rotForce < 0 {
			rotForce = models.Degree(-1 * i.conf.Mecha.RotationForceMaxRange)
		} else if rotForce > 0 {
			rotForce = models.Degree(i.conf.Mecha.RotationForceMaxRange)
		}
	}

	force.AddRot(rotForce)

	// Biar gak overshoot
	if models.Degree(math.Abs(float64(rotForce))) < models.Degree(i.conf.CommandParameter.LookatToleranceDeg) {
		force.ClearRot()
	}
}

func (i LookWRotCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
