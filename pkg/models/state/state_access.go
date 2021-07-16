package state

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type StateAccess struct {
	config       *configuration.FreezeConfig
	lock         sync.Mutex
	myState      RobotState
	stateChecker *time.Ticker
}

func CreateStateAccess(conf *configuration.FreezeConfig) *StateAccess {
	state := &StateAccess{config: conf}
	state.UpdateMyName(conf.Robot.Name)
	state.UpdateMyColor(string(conf.Robot.Color))
	return state
}

// SetMyState digunakan untuk mengubah state diri concurrent safe
// Gunakan fungsi ini untuk menset myState
func (s *StateAccess) SetState(r RobotState) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.myState = r
}

// Fungsi ini berfungsi untuk membaca state diri concurrent safe
func (s *StateAccess) GetState() RobotState {
	// TODO: Do I need to lock for reading?
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.myState
}

func (s *StateAccess) GetStateBytes() ([]byte, error) {
	state := s.GetState()
	b, err := json.Marshal(state)
	if err != nil {
		return nil, errors.New("gagal get state json")
	}
	return b, nil
}

func (s *StateAccess) GetStateJson() (string, error) {
	b, err := s.GetStateBytes()
	if err != nil {
		return "", err
	}
	jsonMsg := string(b)
	return jsonMsg, nil
}

func (s *StateAccess) StopWatcher() {
	s.stateChecker.Stop()
}

func (s *StateAccess) StartWatcher(config *configuration.FreezeConfig) {
	s.stateChecker = time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			<-s.stateChecker.C

			s.lock.Lock()

			// Ball Transform Expiration
			if time.Since(s.myState.BallTransformLastUpdate) > config.Expiration.BallExpiration {
				s.myState.BallTransformExpired = true
			}

			// My Transform Expiration
			if time.Since(s.myState.MyTransformLastUpdate) > config.Expiration.MyExpiration {
				s.myState.MyTransformExpired = true
			}

			// FGP Expiration
			if time.Since(s.myState.FriendGoalPostTransformLastUpdate) > config.Expiration.MyExpiration {
				s.myState.FriendGoalPostTransformExpired = true
			}

			// EGP Expiration
			if time.Since(s.myState.EnemyGoalPostTransformLastUpdate) > config.Expiration.MyExpiration {
				s.myState.EnemyGoalPostTransformExpired = true
			}

			// GutToBrain Expiration
			if time.Since(s.myState.GutToBrainLastUpdate) > config.Expiration.MyExpiration {
				s.myState.GutToBrainExpired = true
			}

			// TODO: Friend, Enemy, Obstacle

			s.lock.Unlock()
		}
	}()
}
