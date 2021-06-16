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

	WADirection WasdDirection = "WA"
	AWDirection WasdDirection = "AW"

	WDDirection WasdDirection = "WD"
	DWDirection WasdDirection = "DW"

	SDDirection WasdDirection = "SD"
	DSDirection WasdDirection = "DS"

	ASDirection WasdDirection = "AS"
	SADirection WasdDirection = "SA"
)

type WasdCommand struct {
	Direction WasdDirection
}

var (
	acceptedDir []WasdDirection = []WasdDirection{
		WDirection,
		ADirection,
		SDirection,
		DDirection,

		WADirection,
		AWDirection,

		WDDirection,
		DWDirection,

		SDDirection,
		DSDirection,

		ASDirection,
		SADirection,
	}
)

func ParseWasdCommand(cmd string) (bool, CommandInterface) {
	dir := strings.ToUpper(strings.TrimSpace(cmd))

	for _, v := range acceptedDir {
		if dir == string(v) {
			return true, WasdCommand{
				Direction: WasdDirection(dir),
			}
		}
	}

	return false, nil
}

func (i WasdCommand) GetName() string {
	return "WASD:" + string(i.Direction)
}

func (i WasdCommand) Tick() {

}
