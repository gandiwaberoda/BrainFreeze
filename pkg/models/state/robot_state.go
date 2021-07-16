package state

import (
	"time"

	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/gutmodel"
)

type RobotState struct {
	MyName           string
	CurrentObjective string
	MyColor          string

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

	MagentaTransform           models.Transform
	MagentaTransformLastUpdate time.Time
	MagentaTransformExpired    bool

	CyanTransform           models.Transform
	CyanTransformLastUpdate time.Time
	CyanTransformExpired    bool

	FriendGoalPostTransform           models.Transform
	FriendGoalPostTransformLastUpdate time.Time
	FriendGoalPostTransformExpired    bool

	EnemyTransform           []models.Transform
	EnemyTransformLastUpdate time.Time
	EnemyTransformExpired    bool

	EnemyGoalPostTransform           models.Transform
	EnemyGoalPostTransformLastUpdate time.Time
	EnemyGoalPostTransformExpired    bool

	ObstacleTransform                []models.Transform
	ObstacleTransformLastUpdate      time.Time
	ObstacleGoalPostTransformExpired bool

	CircularFieldLine []float64 // Untuk line follower

	Register           registerState // Untuk komunikasi antar robot, dipake untuk Wait fulfillment
	RegisterLastUpdate time.Time
	RegisterExpired    bool

	gameState           GameState // Private, kalau dibuat public nanti overwhelming telemetry yang dikirim
	GameStateLastUpdate time.Time
	GameStateExpired    bool
}
