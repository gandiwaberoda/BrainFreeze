package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
)

func (s *StateAccess) UpdateGutToBrain(gtb gutmodel.GutToBrain) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.GutToBrain = gtb
	s.myState.GutToBrainLastUpdate = time.Now()
	s.myState.GutToBrainExpired = false
}
