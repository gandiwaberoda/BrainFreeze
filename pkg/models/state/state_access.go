package state

import (
	"sync"
	"time"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

var (
	lock         sync.Mutex
	myState      RobotState
	stateChecker *time.Ticker
)

// SetMyState digunakan untuk mengubah state diri concurrent safe
// Gunakan fungsi ini untuk menset myState
func SetState(s RobotState) {
	lock.Lock()
	defer lock.Unlock()

	myState = s
}

// Fungsi ini berfungsi untuk membaca state diri concurrent safe
func GetState() RobotState {
	// TODO: Do I need to lock for reading?
	lock.Lock()
	defer lock.Unlock()

	return myState
}

func StopWatcher() {
	stateChecker.Stop()
}

func StartWatcher(config *configuration.FreezeConfig) {
	stateChecker = time.NewTicker(100 * time.Millisecond)

	go func() {
		for {
			<-stateChecker.C

			lock.Lock()

			// Ball Transform Expiration
			if time.Since(myState.BallTransformLastUpdate) > config.Expiration.BallExpiration {
				myState.BallTransformExpired = true
			}

			// My Transform Expiration
			if time.Since(myState.MyTransformLastUpdate) > config.Expiration.MyExpiration {
				myState.MyTransformExpired = true
			}

			// TODO: Friend, EGP, FGP, Enemy

			lock.Unlock()
		}
	}()
}
