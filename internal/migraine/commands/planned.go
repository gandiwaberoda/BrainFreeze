package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PlannedCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	sequence    SequenceCommand
}

func ParsePlannedCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "PLANNED") {
		return false, nil, nil
	}

	seq, err := ParseSequenceCommand("PLANNED", cmd, conf, curstate)
	parsed := PlannedCommand{
		sequence:    seq,
		fulfillment: fulfillments.DefaultComplexFulfillment(),
	}
	if err != nil {
		return true, nil, err
	}

	// Lakukan cek semua subcmd valid
	if err := ValidateSubcmds(seq); err != nil {
		return true, nil, err
	}

	return true, &parsed, nil
}

func (i PlannedCommand) GetName() string {
	if i.sequence.current_obj != nil {
		return "PLANNED [" + i.sequence.current_obj.GetName() + "]"
	} else {
		return "PLANNED [initializing]"
	}
}

func (i *PlannedCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	finished := i.sequence.Tick(force, state)
	if finished {
		i.fulfillment.(*fulfillments.ComplexFuilfillment).Fulfilled()
	}
}

func (i PlannedCommand) ShouldClear() bool {
	return i.fulfillment.ShouldClear()
}

func (i PlannedCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
