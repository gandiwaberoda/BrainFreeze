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
	LastToGut        string

	FpsHsv int

	MyTransform           models.Transform
	MyTransformLastUpdate time.Time
	MyTransformExpired    bool `default:"true"`

	BallTransform           models.Transform
	BallTransformLastUpdate time.Time
	BallTransformExpired    bool `default:"true"`

	GutToBrain           gutmodel.GutToBrain
	GutToBrainLastUpdate time.Time
	GutToBrainExpired    bool `default:"true"`

	Araya           gutmodel.Araya
	ArayaLastUpdate time.Time
	ArayaExpired    bool `default:"true"`

	Straight           []models.StraightDetectionObj
	StraightLastUpdate time.Time
	StraightExpired    bool `default:"true"`

	MagentaTransform           models.Transform
	MagentaTransformLastUpdate time.Time
	MagentaTransformExpired    bool `default:"true"`

	CyanTransform           models.Transform
	CyanTransformLastUpdate time.Time
	CyanTransformExpired    bool `default:"true"`

	FriendGoalPostTransform           models.Transform
	FriendGoalPostTransformLastUpdate time.Time
	FriendGoalPostTransformExpired    bool `default:"true"`

	EnemyTransform           []models.Transform
	EnemyTransformLastUpdate time.Time
	EnemyTransformExpired    bool `default:"true"`

	EnemyGoalPostTransform           models.Transform
	EnemyGoalPostTransformLastUpdate time.Time
	EnemyGoalPostTransformExpired    bool `default:"true"`

	CircularFieldLine []float64 // Untuk line follower

	Register           registerState // Untuk komunikasi antar robot, dipake untuk Wait fulfillment
	RegisterLastUpdate time.Time
	RegisterExpired    bool `default:"true"`

	gameState           GameState // Private, kalau dibuat public nanti overwhelming telemetry yang dikirim
	GameStateLastUpdate time.Time
	GameStateExpired    bool `default:"true"`

	ObstacleTransform           []models.Transform
	ObstacleTransformLastUpdate time.Time
	ObstacleTransformExpired    bool `default:"true"`
}
