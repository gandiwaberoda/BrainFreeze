package diagnostic

import (
	"fmt"
	"log"
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
	"harianugrah.com/brainfreeze/pkg/telepathy"
)

type Telemetry struct {
	isRunning bool
	state     *state.StateAccess
	ticker    *time.Ticker
	tele      telepathy.Telepathy
	config    *configuration.FreezeConfig
}

func CreateNewTelemetry(telepathy telepathy.Telepathy, config *configuration.FreezeConfig, state *state.StateAccess) *Telemetry {
	return &Telemetry{
		isRunning: false,
		tele:      telepathy,
		config:    config,
		state:     state,
	}
}

func worker(t *Telemetry) {
	t.ticker = time.NewTicker(time.Second / t.config.Diagnostic.TelemetryHz)

	for {
		<-t.ticker.C
		tele := t.tele

		json, err := t.state.GetStateJson()
		if err != nil {
			log.Fatalln("Idk what's happening")
		}

		intercom := models.Intercom{
			Kind:     models.TELEMETRY,
			Receiver: models.ALL,
			Content:  json,
		}
		intercomMsg, intercomErr := intercom.AsJson()
		if intercomErr != nil {
			fmt.Println("Gagal jsonify intercomErr", intercomErr.Error())
		}

		_, errTele := tele.Send(intercomMsg)
		if errTele != nil {
			fmt.Println("Failed sending", errTele)
		}

	}
}

func (t *Telemetry) Start() (bool, error) {
	go worker(t)
	return true, nil
}

func (t *Telemetry) Stop() (bool, error) {
	t.ticker.Stop()
	return true, nil
}
