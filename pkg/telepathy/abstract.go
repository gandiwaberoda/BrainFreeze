package telepathy

type Telepathy interface {
	Start() (bool, error)
	Stop() (bool, error)
	Send(string) (bool, error)
	RegisterHandler(func(string))
}
