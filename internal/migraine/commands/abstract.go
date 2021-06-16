package commands

import (
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type CommandInterface interface {
	GetName() string
	Tick(*models.Force, *state.StateAccess)
}
