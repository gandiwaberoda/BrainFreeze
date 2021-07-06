package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateMyName(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.MyName = name
}

func (s *StateAccess) UpdateMyColor(col string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.MyColor = col
}

func (s *StateAccess) UpdateMyTransform(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.MyTransform = t
	s.myState.MyTransformLastUpdate = time.Now()
	s.myState.MyTransformExpired = false
}
