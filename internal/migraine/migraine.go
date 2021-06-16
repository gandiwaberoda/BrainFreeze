package migraine

import (
	"fmt"
	"strings"
	"time"

	"harianugrah.com/brainfreeze/internal/gut"
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type Migraine struct {
	config           *configuration.FreezeConfig
	gut              gut.GutInterface
	ticker           *time.Ticker
	IsRunning        bool
	CurrentObjective commands.CommandInterface
	state            *state.StateAccess
}

func worker(m *Migraine) {
	m.ticker = time.NewTicker(time.Second / m.config.Diagnostic.TelemetryHz)

	for {
		<-m.ticker.C

		if m.CurrentObjective == nil {
			continue
		}

		fmt.Println("Current objective:", m.CurrentObjective.GetName())
		force := models.Force{}
		m.CurrentObjective.Tick(&force, m.state)

		m.gut.Send(force.AsGutCommandString())
	}
}

func CreateMigraine(conf *configuration.FreezeConfig, _gut gut.GutInterface, state *state.StateAccess) *Migraine {
	return &Migraine{config: conf, gut: _gut, state: state}
}

func (m *Migraine) Start() {
	m.Idle()
	go worker(m)
}

func (m *Migraine) Stop() {

}

func (m *Migraine) ReplaceObjective(cmd commands.CommandInterface) {
	m.CurrentObjective = cmd

	str_obj := m.CurrentObjective.GetName()
	m.state.UpdateCurrentObjective(str_obj)
}

func (m *Migraine) AddCommand(intercom models.Intercom) {
	shouldListen := amIReceiver(intercom, m)
	if !shouldListen {
		fmt.Println("I am not a receiver for the command")
		return
	}

	cmd := WhichCommand(intercom, m.config)

	if cmd != nil {
		m.ReplaceObjective(cmd)
	} else {
		fmt.Println("No handler for command")
	}
}

//HELPER
func amIReceiver(intercom models.Intercom, m *Migraine) bool {
	_amIReceiver := false

	myReceiverTag := []string{string(models.ALL), string(m.config.Robot.Name), string(m.config.Robot.Role)}

	for _, v := range myReceiverTag {
		// Case insensitive
		if strings.EqualFold(string(intercom.Receiver), v) {
			_amIReceiver = true
		}
	}
	return _amIReceiver
}

// =========== Basic Command Shorthand ==========
func (m *Migraine) Idle() {
	m.ReplaceObjective(commands.IdleCommand{})
}
