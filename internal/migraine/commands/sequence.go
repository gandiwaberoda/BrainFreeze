package commands

import (
	"errors"
	"fmt"

	"harianugrah.com/brainfreeze/internal/migraine/helper"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type SequenceCommand struct {
	subcommands_str []string // Sudah di ubah spasi menjadi / juga, delimeternya ;
	current_obj     CommandInterface
	conf            *configuration.FreezeConfig
	state           *state.StateAccess
	cmd_name        string // Dipake untuk tampil di CrossRoad
}

func ValidateSubcmds(copied SequenceCommand) error {
	// Cek Apakah sub cmd valid semua
	for _, sub := range copied.subcommands_str {
		_subcmd, _err := WhichCommand(sub, copied.conf, copied.state)
		if _subcmd == nil {
			return errors.New(fmt.Sprint("Sequence subcommand not found:", sub))
		}
		if _err != nil {
			return errors.New(fmt.Sprint("Sequence subcommand are't valid:", _err))
		}
	}

	return nil
}

// Ini tidak dipakai secara langsung, tapi di implement oleh command lain kayak HANDLING sama PLANNED
func ParseSequenceCommand(cmd_name string, cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (SequenceCommand, error) {
	parsed := SequenceCommand{}
	parsed.cmd_name = cmd_name
	parsed.conf = conf
	parsed.state = curstate
	parsed.subcommands_str = cmd.Parameter

	// Lakukan cek semua subcmd valid
	if err := ValidateSubcmds(parsed); err != nil {
		return parsed, err
	}

	return parsed, nil
}

func (i *SequenceCommand) NextObjective() (finished bool) {
	if len(i.subcommands_str) == 0 {
		// Sudah command terakhir
		return true
	}
	// nextup := strings.TrimSpace(i.subcommands_str[0])
	nextup, err := bfvid.ParseCommandSPOK(i.subcommands_str[0])
	if err != nil {
		panic("it should be handled on parsing")
	}
	// fmt.Println("Next obj", nextup)
	i.subcommands_str = removeIndex(i.subcommands_str, 0)

	if nextup.Receiver != "" {
		if !helper.AmIReceiver(nextup.Receiver, i.conf) {
			fmt.Println("Skipped one command (", nextup, ") for [", nextup.Receiver, "] as it is not me")
			return i.NextObjective()
		}
	}

	nextcmd, err := WhichCommand(nextup.Raw, i.conf, i.state)
	if nextcmd == nil {
		panic("INVALID SUBCMD: " + nextup.Raw) // FIXME: Handle pas parse pertama
	}

	if err != nil {
		panic("ERROR Parsing Subcmd: " + nextup.Raw) // FIXME: Handle pas parse pertama
	}

	i.current_obj = nextcmd
	str_obj := i.cmd_name + " [" + i.current_obj.GetName() + " -> " + i.current_obj.GetFulfillment().AsString() + "]"
	i.state.UpdateCurrentObjective(str_obj)
	return false
}

func (i *SequenceCommand) Tick(force *models.Force, state *state.StateAccess) (finished bool) {
	if i.current_obj == nil {
		fmt.Println("nilll")
		return i.NextObjective()
	}

	i.current_obj.Tick(force, state)
	if i.current_obj.GetFulfillment().ShouldClear() {
		fmt.Println("seq next")
		if i.NextObjective() {
			return true
		}
	}

	return false
}

// func (i PlannedCommand) ShouldClear() bool {
// 	return i.fulfillment.ShouldClear()
// }

// func (i PlannedCommand) GetFulfillment() fulfillments.FulfillmentInterface {
// 	return i.fulfillment
// }

// Helper
func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
