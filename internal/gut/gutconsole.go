package gut

import (
	"fmt"
)

type GutConsole struct {
}

func CreateGutConsole() *GutConsole {
	return &GutConsole{}
}

func (g *GutConsole) Start() (bool, error) {
	return true, nil
}

func (g *GutConsole) Stop() (bool, error) {
	return true, nil
}

func (g *GutConsole) Send(msg string) (bool, error) {
	fmt.Println("TO GUT:", msg)
	return true, nil
}

func (g *GutConsole) RegisterHandler(handler func(string)) {

}
