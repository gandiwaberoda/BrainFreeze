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
type MoveWXCommand struct {
	StartTransform models.Transform
	SpeedX         int
	conf           *configuration.FreezeConfig
	fulfillment    fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseMoveWXCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "MOVEWX") {
		return false, nil, nil
	}

	targetMX := 0
	if len(cmd.Parameter) == 1 && cmd.Parameter[0] != "" {
		d, err := strconv.Atoi(cmd.Parameter[0])
		if err != nil {
			return true, nil, errors.New("move w x failed convert to int")
		}
		targetMX = d
	} else {
		return true, nil, errors.New("move w x command require 1 parameter")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultDeltaposFulfillment(models.Centimeter(math.Abs(float64(targetMX))), curstate)
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := MoveWXCommand{
		StartTransform: curstate.GetState().MyTransform,
		SpeedX:         targetMX,
		conf:           conf,
		fulfillment:    parsedFulfilment,
	}

	return true, &parsed, nil
}

func (i MoveWXCommand) GetName() string {
	return fmt.Sprint("MOVE WORLD X: ", i.SpeedX)
}

func (i *MoveWXCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	TockGoto(float64(state.GetState().MyTransform.WorldXcm+models.Centimeter(i.SpeedX)), float64(i.StartTransform.WorldYcm), i.conf, force, state)
	force.ClampMinXY(*i.conf)
}

func (i MoveWXCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
