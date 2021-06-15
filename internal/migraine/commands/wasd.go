package commands

import (
	"strings"
)

type WasdDirection string

const (
	WDirection WasdDirection = "W"
	ADirection WasdDirection = "A"
	SDirection WasdDirection = "S"
	DDirection WasdDirection = "D"
)

type WasdCommand struct {
	Direction WasdDirection
}

func ParseWasdCommand(cmd string) (bool, CommandInterface) {
	fLetter := strings.ToUpper(cmd[:1])
	if fLetter == "W" || fLetter == "A" || fLetter == "S" || fLetter == "D" {
		return true, WasdCommand{
			Direction: WasdDirection(fLetter),
		}
	}

	return false, nil
}

func (i WasdCommand) GetName() string {
	return "WASD:" + string(i.Direction)
}

func (i WasdCommand) Tick() {

}
