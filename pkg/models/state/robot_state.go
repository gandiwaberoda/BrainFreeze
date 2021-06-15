package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/gut"
)

type RobotState struct {
	MyName           string
	CurrentObjective string

	MyTransform           models.Transform
	MyTransformLastUpdate time.Time
	MyTransformExpired    bool

	BallTransform           models.Transform
	BallTransformLastUpdate time.Time
	BallTransformExpired    bool

	GutToBrain                    gut.GutToBrain
	GutToBrainTransformLastUpdate time.Time
	GutToBrainTransformExpired    bool

	FriendTransform           []models.Transform
	FriendTransformLastUpdate time.Time

	FriendGoalPostTransform           []models.Transform
	FriendGoalPostTransformLastUpdate time.Time

	EnemyTransform           []models.Transform
	EnemyTransformLastUpdate time.Time

	EnemyGoalPostTransform           []models.Transform
	EnemyGoalPostTransformLastUpdate time.Time

	ObstacleTransform           []models.Transform
	ObstacleTransformLastUpdate time.Time
}
