package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateStraight(st []models.StraightDetectionObj) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.Straight = st
	s.myState.StraightLastUpdate = time.Now()
	s.myState.StraightExpired = false
}
