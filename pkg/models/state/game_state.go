package state

import (
	"encoding/json"
	"errors"
	"time"
)

type GameState struct {
	RobotStates []RobotState
}

func (s *StateAccess) UpdateGameState(gs GameState) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState.gameState = gs
	s.myState.GameStateLastUpdate = time.Now()
	s.myState.GameStateExpired = false
}

func (s GameState) AsBytes() ([]byte, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, errors.New("gagal get bytes")
	}
	return b, nil
}

func (s GameState) AsJson() (string, error) {
	b, err := s.AsBytes()
	if err != nil {
		return "", err
	}
	jsonMsg := string(b)
	return jsonMsg, nil
}
