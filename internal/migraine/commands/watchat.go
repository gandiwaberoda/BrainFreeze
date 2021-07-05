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

type WatchatCommand struct {
	Target      string
	conf        *configuration.FreezeConfig
	fulfillment fulfillments.FulfillmentInterface
}

// WasdCommand memiliki fulfillment default yaitu DefaultDurationFulfillment
func ParseWatchatCommand(intercom models.Intercom, cmd string, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface) {
	if len(cmd) < 7 {
		return false, nil
	}

	if !strings.EqualFold(cmd[:7], "WATCHAT") {
		return false, nil
	}

	re, _ := regexp.Compile(`\(([A-Za-z0-9]+)\)`)
	foundParam := re.FindString(cmd)
	foundParam = strings.ReplaceAll(foundParam, "(", "")
	foundParam = strings.ReplaceAll(foundParam, ")", "")

	target := "BALL"
	if foundParam != "" {
		fmt.Println(foundParam)
		target = foundParam
	}

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		return false, nil
	}

	parsed := WatchatCommand{
		Target:      target,
		conf:        conf,
		fulfillment: fulfillments.DefaultHoldFulfillment(),
	}

	return true, &parsed
}

func (i WatchatCommand) GetName() string {
	return "WATCHAT:" + string(i.Target)
}

func (i *WatchatCommand) Tick(force *models.Force, state *state.StateAccess) {
	_, target := state.GetTransformByKey(i.Target)

	TockLookat(target, *i.conf, force, state)

	i.fulfillment.Tick()
}

func (i WatchatCommand) GetFulfillment() fulfillments.FulfillmentInterface {
	return i.fulfillment
}
