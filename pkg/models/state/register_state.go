package state

import (
	"time"
)

type registerState struct {
	ReadyReceive float64
}

type RegisterKey string

const (
	READY_RECEIVED RegisterKey = "ReadyReceived"
)

func (s *StateAccess) UpdateRegisterState(gs registerState) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.Register = gs
	s.myState.RegisterLastUpdate = time.Now()
	s.myState.RegisterExpired = false
}

func NewRegister() registerState {
	return registerState{
		ReadyReceive: 0.0,
	}
}
