package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type AvoidanceCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	sequence    SequenceCommand
}

func ParseAvoidanceCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "AVOIDANCE") {
		return false, nil, nil
	}

	seq, err := ParseSequenceCommand("AVOIDANCE", cmd, conf, curstate)
	parsed := AvoidanceCommand{
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

func (i AvoidanceCommand) GetName() string {
	if i.sequence.current_obj != nil {
		return "AVOIDANCE (" + i.sequence.current_obj.GetName() + ") [" + i.sequence.current_obj.GetName() + "]"
	} else {
		return "AVOIDANCE [initializing]"
	}
}

func (i *AvoidanceCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()
	finished := i.sequence.Tick(force, state)
	force.EnableHandling()
	if finished {
		i.fulfillment.(*fulfillments.ComplexFuilfillment).Fulfilled()
	}
}

func (i AvoidanceCommand) ShouldClear() bool {
	return i.fulfillment.ShouldClear()
}

func (i AvoidanceCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
