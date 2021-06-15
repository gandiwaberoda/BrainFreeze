package migraine

import (
	"fmt"
	"time"

	"harianugrah.com/brainfreeze/internal/gut"
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type Migraine struct {
	config           *configuration.FreezeConfig
	gut              gut.GutInterface
	ticker           *time.Ticker
	IsRunning        bool
	CurrentObjective commands.CommandInterface
}

func worker(m *Migraine) {
	m.ticker = time.NewTicker(time.Second / m.config.Diagnostic.TelemetryHz)

	for {
		<-m.ticker.C
		if m.CurrentObjective == nil {
			continue
		}

		fmt.Println("Current objective:", m.CurrentObjective.GetName())
		m.gut.Send("Apalah")
	}
}

func CreateMigraine(conf *configuration.FreezeConfig, _gut gut.GutInterface) *Migraine {
	return &Migraine{config: conf, gut: _gut}
}

func (m *Migraine) Start() {
	m.Idle()
	go worker(m)
}

func (m *Migraine) Stop() {

}

func (m *Migraine) Idle() {
	m.CurrentObjective = commands.IdleCommand{}
}

func (m *Migraine) AddCommand(intercom models.Intercom) {
	cmd := ParseCommand(intercom)
	if cmd != nil {
		m.CurrentObjective = cmd
	} else {
		fmt.Println("No handler for command")
	}
}
