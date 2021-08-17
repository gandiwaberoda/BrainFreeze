package state

func (s *StateAccess) UpdateFpsHsv(fps int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.FpsHsv = fps
}

func (s *StateAccess) UpdateToGutCmd(str string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.LastToGut = str
}
