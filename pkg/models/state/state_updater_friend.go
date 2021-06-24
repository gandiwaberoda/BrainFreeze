package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func (s *StateAccess) UpdateMagentaTransform(t models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.MagentaTransform = t
	s.myState.MagentaTransformLastUpdate = time.Now()
	s.myState.MagentaTransformExpired = false
}
