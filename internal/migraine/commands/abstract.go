package commands

type CommandInterface interface {
	GetName() string
	Tick()
}
