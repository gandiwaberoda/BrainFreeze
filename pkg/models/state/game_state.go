package state

import (
	"encoding/json"
	"errors"
	"strings"
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

func (s *StateAccess) GetOtherRegisterByIdentifier(color_or_name string, key RegisterKey) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.myState.GameStateExpired {
		return 0, errors.New("game state is expired")
	}

	for _, v := range s.myState.gameState.RobotStates {
		if strings.EqualFold(v.MyColor, color_or_name) || strings.EqualFold(v.MyName, color_or_name) {
			if v.RegisterExpired {
				return 0, errors.New("register of the friend robot is expired")
			}

			if key == READY_RECEIVED {
				return v.Register.ReadyReceive, nil
			}
			// TODO: Tambahkan register access yang lain disini
		}
	}

	return 0, errors.New("color/name (Identifier) not found")
}
