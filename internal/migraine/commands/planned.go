package commands

import (
	"fmt"
	"regexp"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type PlannedCommand struct {
	fulfillment        fulfillments.FulfillmentInterface
	subcommand_raw_str string
	subcommands_str    []string // Sudah di ubah spasi menjadi / juga, delimeternya ;
	current_obj        CommandInterface
	intercom           models.Intercom
	conf               *configuration.FreezeConfig
	// shouldClear        bool
	state *state.StateAccess
}

func ParsePlannedCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 7 {
		return false, nil
	}

	if !strings.EqualFold(cmd[:7], "PLANNED") {
		return false, nil
	}

	parsed := PlannedCommand{}
	parsed.subcommand_raw_str = cmd
	parsed.intercom = intercom
	parsed.conf = conf
	parsed.fulfillment = fulfillments.DefaultComplexFulfillment()
	parsed.state = curstate

	re, _ := regexp.Compile(`\((.+)\)`)
	foundParam := re.FindStringSubmatch(cmd)
	if len(foundParam) < 1 {
		return false, nil
	}

	subcmd := foundParam[0]
	subcmd = subcmd[1 : len(subcmd)-1]
	subcmd = strings.ReplaceAll(subcmd, "@", "/")

	fmt.Println("FP:", foundParam, "--", subcmd)

	parsed.subcommands_str = strings.Split(subcmd, ";")
	// parsed.NextObjective()

	return true, &parsed
}

func (i *PlannedCommand) NextObjective() (finished bool) {
	if len(i.subcommands_str) == 0 {
		// Sudah command terakhir
		return true
	}
	nextup := i.subcommands_str[0]
	fmt.Println("Next obj", nextup)
	i.subcommands_str = removeIndex(i.subcommands_str, 0)

	splitted := strings.Split(string(nextup), ";")

	// inkom_content := string(i.intercom.Kind) + "/"
	inkom_content := ""

	if len(splitted) == 0 {
		panic("...wthat")
	} else if len(splitted) == 1 {
		inkom_content += strings.TrimSpace(splitted[0])
	} else if len(splitted) == 2 {
		inkom_content += strings.TrimSpace(splitted[0]) + "/" + strings.TrimSpace(splitted[1])
	}

	inkom := models.Intercom{
		Kind:     i.intercom.Kind,
		Receiver: i.intercom.Receiver,
		Content:  inkom_content,
	}

	nextcmd := WhichCommand(inkom, i.conf, i.state)
	if nextcmd == nil {
		fmt.Println("INVALID SUBCMD: " + nextup)
		return i.NextObjective()
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

	i.fulfillment.Tick()
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
