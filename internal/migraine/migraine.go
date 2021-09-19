package migraine

import (
	"fmt"
	"time"

	"harianugrah.com/brainfreeze/internal/gut"
	"harianugrah.com/brainfreeze/internal/migraine/commands"
	"harianugrah.com/brainfreeze/internal/migraine/helper"
	"harianugrah.com/brainfreeze/pkg/bfvid"
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

		// fmt.Println("Current objective:", m.CurrentObjective.GetName())

		force := models.Force{}

		register := state.NewRegister()
		m.state.UpdateRegisterState(register)

		m.CurrentObjective.Tick(&force, m.state)

		m.gut.Send(force.AsGutCommandString())

		if m.CurrentObjective.GetFulfillment().ShouldClear() {
			fmt.Println("Fulfilled, now idling")
			m.ReplaceObjective(commands.DefaultIdleCommand())
		}

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
	// TODO
	// Agar saat software ditutup, gak nyantol terus bergerak
	idleForce := models.Force{}
	m.gut.Send(idleForce.AsGutCommandString())
	fmt.Println("whhhhh")
}

func (m *Migraine) ReplaceObjective(cmd commands.CommandInterface) {
	fmt.Println("Replace Objective:", cmd.GetName())
	m.CurrentObjective = cmd

	str_obj := m.CurrentObjective.GetName() + " -> " + m.CurrentObjective.GetFulfillment().AsString()
	m.state.UpdateCurrentObjective(str_obj)
}

func (m *Migraine) AddCommand(fullbfvid string) {
	fmt.Println("ADD COMMAND:", fullbfvid)

	parsed, err := bfvid.ParseCommandSPOK(fullbfvid)
	if err != nil {
		fmt.Println("failed to parse command:", err)
		return
	}

	if parsed.Receiver != "" {
		shouldListen := helper.AmIReceiver(string(parsed.Receiver), m.config)
		if !shouldListen {
			fmt.Println("I am not a receiver for the command")
			return
		}
	}

	// if strings.EqualFold(parsed.Verb, "FWD") {
	// 	// TODO: Command khusus untuk mengforward data serial
	// 	// re, _ := regexp.Compile(`\((.+)\)`)
	// 	// foundParam := re.FindString(intercom.Content)
	// 	// foundParam = strings.ReplaceAll(foundParam, "(", "")
	// 	// foundParam = strings.ReplaceAll(foundParam, ")", "")
	// 	m.gut.Send(parsed.ParameterStr)
	// 	m.Idle()
	// 	return
	// }

	cmd, err := commands.WhichCommand(fullbfvid, m.config, m.state)
	if err != nil {
		fmt.Println("MIGRAINE ERROR:", err)
		return
	}

	if cmd != nil {
		fmt.Println("New Objective Understood")
		m.ReplaceObjective(cmd)
	} else {
		fmt.Println("No handler for command")
		// TODO: Kirim log ke basestation
	}
}

// =========== Basic Command Shorthand ==========
func (m *Migraine) Idle() {
	m.ReplaceObjective(commands.DefaultIdleCommand())
}
