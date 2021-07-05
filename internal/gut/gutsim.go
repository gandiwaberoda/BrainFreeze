package gut

import (
	"harianugrah.com/brainfreeze/internal/simpserv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type GutSim struct {
	Gut
	conf   *configuration.FreezeConfig
	simpWs *simpserv.SimpWs
}

func CreateGutSim(conf *configuration.FreezeConfig, simpws *simpserv.SimpWs) *GutSim {
	return &GutSim{
		conf:   conf,
		simpWs: simpws,
	}
}

func (g *GutSim) Start() (bool, error) {
	return true, nil
}

func (g *GutSim) Stop() (bool, error) {
	return true, nil
}

func (g *GutSim) Send(msg string) (bool, error) {
	g.simpWs.Broadcast(msg)
	return true, nil
}

func (g *GutSim) RegisterHandler(handler func(string)) {
	panic("Use simpserver on msg received")
}
