package state

func (s *StateAccess) UpdateCircularFieldLine(t []float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.CircularFieldLine = t
}
