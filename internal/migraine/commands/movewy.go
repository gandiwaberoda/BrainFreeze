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

// Move World Y axis
type MoveWYCommand struct {
	StartTransform models.Transform
	SpeedY         int
	conf           *configuration.FreezeConfig
	fulfillment    fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseMoveWYCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "MOVEWY") {
		return false, nil, nil
	}

	targetMY := 0
	if len(cmd.Parameter) == 1 && cmd.Parameter[0] != "" {
		d, err := strconv.Atoi(cmd.Parameter[0])
		if err != nil {
			return true, nil, errors.New("move w y failed convert to int")
		}
		targetMY = d
	} else {
		return true, nil, errors.New("move w y command require 1 parameter")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultDeltaposFulfillment(models.Centimeter(math.Abs(float64(targetMY))), curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := MoveWYCommand{
		StartTransform: curstate.GetState().MyTransform,
		SpeedY:         targetMY,
		conf:           conf,
		fulfillment:    parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i MoveWYCommand) GetName() string {
	return fmt.Sprint("MOVE WORLD Y: ", i.SpeedY)
}

func (i *MoveWYCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	TockGoto(float64(i.StartTransform.WorldXcm), float64(state.GetState().MyTransform.WorldYcm+models.Centimeter(i.SpeedY)), i.conf, force, state)
	force.ClampMinXY(*i.conf)
}

func (i MoveWYCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
