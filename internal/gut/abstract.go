package gut

import (
	"harianugrah.com/brainfreeze/internal/araya"
)

type Gut struct {
	ArayaSens araya.ArayaSerial
	handlers  []func(string)
}

type GutInterface interface {
	Start() (bool, error)
	Stop() (bool, error)
	Send(string) (bool, error)
	RegisterHandler(func(string))
}
