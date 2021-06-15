package commands

import (
	"strings"
)

type IdleCommand struct {
}

func ParseIdleCommand(cmd string) (bool, CommandInterface) {
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

func (i IdleCommand) Tick() {

}
