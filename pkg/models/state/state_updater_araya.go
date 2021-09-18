package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
)

func (s *StateAccess) UpdateAraya(ar gutmodel.Araya) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.Araya = ar
	s.myState.ArayaLastUpdate = time.Now()
	s.myState.ArayaExpired = false
}
