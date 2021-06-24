package gut

type IgnoreConsole struct {
	Gut
}

func CreateIgnoreConsole() *IgnoreConsole {
	return &IgnoreConsole{}
}

func (g *IgnoreConsole) Start() (bool, error) {
	return true, nil
}

func (g *IgnoreConsole) Stop() (bool, error) {
	return true, nil
}

func (g *IgnoreConsole) Send(msg string) (bool, error) {
	return true, nil
}

func (g *IgnoreConsole) RegisterHandler(handler func(string)) {

}
