package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type HandlingCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	sequence    SequenceCommand
}

func ParseHandlingCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "HANDLING") {
		return false, nil, nil
	}

	seq, err := ParseSequenceCommand("HANDLING", cmd, conf, curstate)
	parsed := HandlingCommand{
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

func (i HandlingCommand) GetName() string {
	if i.sequence.current_obj != nil {
		return "HANDLING (" + i.sequence.current_obj.GetName() + ") [" + i.sequence.current_obj.GetName() + "]"
	} else {
		return "HANDLING [initializing]"
	}
}

func (i *HandlingCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	finished := i.sequence.Tick(force, state)
	force.EnableHandling()
	if finished {
		i.fulfillment.(*fulfillments.ComplexFuilfillment).Fulfilled()
	}
}

func (i HandlingCommand) ShouldClear() bool {
	return i.fulfillment.ShouldClear()
}

func (i HandlingCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
