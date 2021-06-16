package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type IdleCommand struct {
}

func ParseIdleCommand(intercom models.Intercom, conf *configuration.FreezeConfig) (bool, CommandInterface) {
	cmd := intercom.Content

	if len(cmd) < 4 {
		return false, nil
	}

	if strings.ToUpper(cmd[:4]) == "IDLE" {
		return true, IdleCommand{}
	}

	return false, nil
}

func (i IdleCommand) GetName() string {
	return "IDLE"
}

func (i IdleCommand) Tick(force *models.Force, state *state.StateAccess) {

}

func (i IdleCommand) ShouldClear() bool {
	return false
}
