package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateBallTransform(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.BallTransform = t
	s.myState.BallTransformLastUpdate = time.Now()
	s.myState.BallTransformExpired = false
}
