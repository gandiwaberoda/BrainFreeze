package diagnostic

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/telepathy"
)

type Telemetry struct {
	isRunning bool
	ticker    *time.Ticker
	tele      *telepathy.Telepathy
	config    *configuration.FreezeConfig
}

func CreateNewTelemetry(telepathy *telepathy.Telepathy, config *configuration.FreezeConfig) *Telemetry {
	return &Telemetry{
		isRunning: false,
		tele:      telepathy,
		config:    config,
	}
}

func worker(t *Telemetry) {
	t.ticker = time.NewTicker(time.Second / t.config.Diagnostic.TelemetryHz)

	for {
		<-t.ticker.C
		tele := *t.tele
		tele.Send(t.)
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
