package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
)

type RobotState struct {
	MyName           string
	CurrentObjective string

	FpsHsv int

	MyTransform           models.Transform
	MyTransformLastUpdate time.Time
	MyTransformExpired    bool

	BallTransform           models.Transform
	BallTransformLastUpdate time.Time
	BallTransformExpired    bool

	GutToBrain           gutmodel.GutToBrain
	GutToBrainLastUpdate time.Time
	GutToBrainExpired    bool

	FriendTransform           []models.Transform
	FriendTransformLastUpdate time.Time

	FriendGoalPostTransform           models.Transform
	FriendGoalPostTransformLastUpdate time.Time

	EnemyTransform           []models.Transform
	EnemyTransformLastUpdate time.Time

	EnemyGoalPostTransform           models.Transform
	EnemyGoalPostTransformLastUpdate time.Time

	ObstacleTransform           []models.Transform
	ObstacleTransformLastUpdate time.Time
}
