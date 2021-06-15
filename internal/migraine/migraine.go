package migraine

import "harianugrah.com/brainfreeze/pkg/models"
import "harianugrah.com/brainfreeze/internal/migraine/commands"

type Migraine struct {
	CurrentObjective commands.Command
}

func worker() {

}

func CreateMigraine() *Migraine {
	return &Migraine{}
}

func (m *Migraine) Idle() {

}

func (m *Migraine) AddCommand(intercom models.Intercom) {

}
