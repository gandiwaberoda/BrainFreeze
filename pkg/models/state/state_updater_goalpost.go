package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateFriendGoalpost(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.FriendGoalPostTransform = t
	s.myState.BallTransformLastUpdate = time.Now()
	s.myState.BallTransformExpired = false
}
