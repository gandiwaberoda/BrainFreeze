package commands

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
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
	conf      *configuration.FreezeConfig
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

func ParseWasdCommand(cmd string, conf *configuration.FreezeConfig) (bool, CommandInterface) {
	dir := strings.ToUpper(strings.TrimSpace(cmd))

	for _, v := range acceptedDir {
		if dir == string(v) {
			return true, &WasdCommand{
				Direction: WasdDirection(dir),
				conf:      conf,
			}
		}
	}

	return false, nil
}

func (i WasdCommand) GetName() string {
	return "WASD:" + string(i.Direction)
}

func TockWasd(dir WasdDirection, conf configuration.FreezeConfig, force *models.Force, state *state.StateAccess) {
	if dir == WDirection {
		// Maju
		force.AddY(float64(conf.Mecha.VerticalForceRange))
	} else if dir == SDirection {
		// Mundur
		force.AddY(-1 * float64(conf.Mecha.VerticalForceRange))
	} else if dir == ADirection {
		// Kiri
		force.AddX(-1 * float64(conf.Mecha.HorizontalForceRange))
	} else if dir == DDirection {
		// Kanan
		force.AddX(float64(conf.Mecha.HorizontalForceRange))
	} else if dir == AWDirection || dir == WADirection {
		// Kiri Depan
		force.AddY(float64(conf.Mecha.VerticalForceRange))
		force.AddX(-1 * float64(conf.Mecha.HorizontalForceRange))
	} else if dir == WDDirection || dir == DWDirection {
		// Kanan Depan
		force.AddY(float64(conf.Mecha.VerticalForceRange))
		force.AddX(float64(conf.Mecha.HorizontalForceRange))
	} else if dir == ASDirection || dir == SADirection {
		// Kiri Belakang
		force.AddY(-1 * float64(conf.Mecha.VerticalForceRange))
		force.AddX(-1 * float64(conf.Mecha.HorizontalForceRange))
	} else if dir == SDDirection || dir == DSDirection {
		// Kanan Belakang
		force.AddY(-1 * float64(conf.Mecha.VerticalForceRange))
		force.AddX(float64(conf.Mecha.HorizontalForceRange))
	} else {
		panic("what the heck happened? " + string(dir))
	}
}

func (i *WasdCommand) Tick(force *models.Force, state *state.StateAccess) {
	TockWasd(i.Direction, *i.conf, force, state)
}
