package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateFriendGoalpostTransform(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.FriendGoalPostTransform = t
	s.myState.FriendGoalPostTransformLastUpdate = time.Now()
	s.myState.FriendGoalPostTransformExpired = false
}

func (s *StateAccess) UpdateEnemyGoalpostTransform(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.EnemyGoalPostTransform = t
	s.myState.EnemyGoalPostTransformLastUpdate = time.Now()
	s.myState.EnemyGoalPostTransformExpired = false
}
