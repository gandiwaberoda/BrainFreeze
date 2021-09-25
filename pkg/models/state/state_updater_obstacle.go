package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateObstaclesTransform(ts []models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.ObstacleTransform = ts
	s.myState.ObstacleTransformLastUpdate = time.Now()
	s.myState.ObstacleTransformExpired = false
}
func (s *StateAccess) UpdateCircularDummy(t []float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.CircularDummy = t
}
