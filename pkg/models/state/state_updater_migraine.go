package state

import ()

func (s *StateAccess) UpdateCurrentObjective(obj string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.CurrentObjective = obj
}
