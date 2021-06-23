package state

import (
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
)

var ACCEPTABLE_TRANSFORM_KEY = map[string]bool{"EGP": true, "FGP": true, "BALL": true}

func GetTransformKeyAcceptable(key string) bool {
	_, found := ACCEPTABLE_TRANSFORM_KEY[key]
	return found
}

// Return found, Transform
func (s *StateAccess) GetTransformByKey(key string) (bool, models.Transform) {
	s.lock.Lock()
	defer s.lock.Unlock()

	target := strings.ToUpper(key)
	if !GetTransformKeyAcceptable(target) {
		return false, models.Transform{}
	}

	switch target {
	case "EGP": // Enemy Goal Post
		return true, s.myState.EnemyGoalPostTransform
	case "FGP": // Enemy Goal Post
		return true, s.myState.FriendGoalPostTransform
	case "BALL": // Enemy Goal Post
		return true, s.myState.BallTransform
	default:
		return false, models.Transform{}
	}

}
