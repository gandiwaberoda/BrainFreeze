package commands

import (
	"errors"
	"fmt"
	"strings"

	"harianugrah.com/brainfreeze/internal/migraine/fulfillments"
	"harianugrah.com/brainfreeze/pkg/bfvid"
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
func ParseWatchatCommand(cmd bfvid.CommandSPOK, conf *configuration.FreezeConfig, curstate *state.StateAccess) (bool, CommandInterface, error) {
	// if len(cmd) < 7 {
	// 	return false, nil
	// }

	// if !strings.EqualFold(cmd[:7], "WATCHAT") {
	// 	return false, nil
	// }
	if !strings.EqualFold(cmd.Verb, "WATCHAT") {
		return false, nil, nil
	}

	// re, _ := regexp.Compile(`\(([A-Za-z0-9]+)\)`)
	// foundParam := re.FindString(cmd)
	// foundParam = strings.ReplaceAll(foundParam, "(", "")
	// foundParam = strings.ReplaceAll(foundParam, ")", "")

	target := "BALL"
	if cmd.Parameter[0] != "" {
		fmt.Println(cmd.Parameter[0])
		target = cmd.Parameter[0]
	}

	isKeyAcceptable := state.GetTransformKeyAcceptable(target)
	if !isKeyAcceptable {
		return true, nil, errors.New("watchat target key not acceptable")
	}

	var parsedFulfilment fulfillments.FulfillmentInterface
	if cmd.Fulfilment == "" {
		parsedFulfilment = fulfillments.DefaultHoldFulfillment()
	} else {
		filment, err := fulfillments.WhichFulfillment(cmd.Raw, conf, curstate)
		if err != nil {
			return true, nil, errors.New(fmt.Sprint("non default fulfilment error:", err))
		}
		parsedFulfilment = filment
	}

	parsed := WatchatCommand{
		Target:      target,
		conf:        conf,
		fulfillment: parsedFulfilment,
	}

	return true, &parsed, nil
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
