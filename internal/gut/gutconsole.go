package gut

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
	return true, nil
}

func (g *GutConsole) RegisterHandler(handler func(string)) {

}
