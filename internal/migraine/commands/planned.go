package commands

import (
	"errors"
	"fmt"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/internal/migraine/helper"
	"harianugrah.com/brainfreeze/pkg/bfvid"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PlannedCommand struct {
	fulfillment fulfillments.FulfillmentInterface
	// subcommand_raw_str string
	subcommands_str []string // Sudah di ubah spasi menjadi / juga, delimeternya ;
	current_obj     CommandInterface
	// intercom           models.Intercom
	conf *configuration.FreezeConfig
	// shouldClear        bool
	state *state.StateAccess
}

func ValidateSubcmds(copied PlannedCommand) error {
	// Cek Apakah sub cmd valid semua
	for _, sub := range copied.subcommands_str {
		_subcmd, _err := WhichCommand(sub, copied.conf, copied.state)
		if _subcmd == nil {
			return errors.New(fmt.Sprint("planned subcommand not found:", sub))
		}
		if _err != nil {
			return errors.New(fmt.Sprint("planned subcommand are't valid:", _err))
		}
	}

	return nil
}

func ParsePlannedCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	if !strings.EqualFold(cmd.Verb, "PLANNED") {
		return false, nil, nil
	}

	parsed := PlannedCommand{}
	parsed.conf = conf
	parsed.fulfillment = fulfillments.DefaultComplexFulfillment()
	parsed.state = curstate
	parsed.subcommands_str = cmd.Parameter

	// Lakukan cek semua subcmd valid
	if err := ValidateSubcmds(parsed); err != nil {
		return true, nil, err
	}

	return true, &parsed, nil
}

func (i *PlannedCommand) NextObjective() (finished bool) {
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

	// splitted := strings.Split(string(nextup), ";")
	// fmt.Println("SPLITTED: ", splitted)
	// if colonIndex := strings.Index(nextup, ":"); colonIndex != -1 {
	// 	// Kalau ada tanda : di subcmd, berarti hanya robot tertentu yang perlu denger
	// 	receiver := strings.TrimSpace(nextup[0:colonIndex])
	// 	nextup = strings.TrimSpace(nextup[colonIndex+1:])

	// }

	// inkom_content := string(i.intercom.Kind) + "/"
	// inkom_content := ""

	// inkom_content += strings.TrimSpace(nextup)

	// inkom := models.Intercom{
	// 	Kind:     i.intercom.Kind,
	// 	Receiver: i.intercom.Receiver,
	// 	Content:  inkom_content,
	// }

	nextcmd, err := WhichCommand(nextup.Raw, i.conf, i.state)
	if nextcmd == nil {
		panic("INVALID SUBCMD: " + nextup.Raw) // FIXME: Handle pas parse pertama
	}

	if err != nil {
		panic("ERROR Parsing Subcmd: " + nextup.Raw) // FIXME: Handle pas parse pertama
	}

	i.current_obj = nextcmd
	str_obj := "PLANNED [" + i.current_obj.GetName() + " -> " + i.current_obj.GetFulfillment().AsString() + "]"
	i.state.UpdateCurrentObjective(str_obj)
	return false
}

func (i PlannedCommand) GetName() string {
	fmt.Println("...")
	if i.current_obj != nil {
		return "PLANNED [" + i.current_obj.GetName() + "]"
	} else {
		return "PLANNED [initializing]"
	}
}

func (i *PlannedCommand) Tick(force *models.Force, state *state.StateAccess) {
	i.fulfillment.Tick()

	if i.current_obj == nil {
		fmt.Println("nilll")
		if i.NextObjective() {
			i.fulfillment.(*fulfillments.ComplexFuilfillment).Fulfilled()
		}
		return
	}

	i.current_obj.Tick(force, state)
	if i.current_obj.GetFulfillment().ShouldClear() {
		fmt.Println("planned next")
		if i.NextObjective() {
			i.fulfillment.(*fulfillments.ComplexFuilfillment).Fulfilled()
		}
	}
}

func (i PlannedCommand) ShouldClear() bool {
	return i.fulfillment.ShouldClear()
}

func (i PlannedCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}

// Helper
func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
