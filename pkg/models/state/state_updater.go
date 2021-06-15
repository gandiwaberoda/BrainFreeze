package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
)

func UpdateMyTransform(t models.Transform) {
	lock.Lock()
	defer lock.Unlock()

	myState.MyTransform = t
	myState.MyTransformLastUpdate = time.Now()
	myState.MyTransformExpired = false
}
