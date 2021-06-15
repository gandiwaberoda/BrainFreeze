package gut

type Gut struct {
	handlers []func(string)
}

type GutInterface interface {
	Start() (bool, error)
	Stop() (bool, error)
	Send(string) (bool, error)
	RegisterHandler(func(string))
}
