package state

import (
	"sort"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models"
)

var ACCEPTABLE_TRANSFORM_KEY = map[string]bool{"EGP": true, "FGP": true, "BALL": true, "MAGENTA": true, "CYAN": true, "CLOBS": true}

func GetTransformKeyAcceptable(key string) bool {
	_, found := ACCEPTABLE_TRANSFORM_KEY[strings.ToUpper(key)]
	return found
}

// Untuk mencari titik yang paling masuk akal menjadi bola, jika diketahui lokasi bola sebelumnya
func SortTransformByRobDistance(other []models.Transform) []models.Transform {
	sort.Slice(other, func(i, j int) bool {
		return other[i].TopRpx < other[j].TopRpx
	})

	return other
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
	case "MAGENTA": // Closest Dummy
		return true, s.myState.MagentaTransform
	case "CYAN": // Closest Dummy
		return true, s.myState.CyanTransform
	case "CLOBS": // Closest Obstacle
		sorted := SortTransformByRobDistance(s.myState.ObstacleTransform)
		return true, sorted[0]
	default:
		return false, models.Transform{}
	}

}
