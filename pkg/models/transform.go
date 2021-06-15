package models

type Centimeter float64
type Degree float64

type Transform struct {
	// cm
	EncXcm Centimeter
	// cm
	EncYcm Centimeter
	// derajat
	EncROT Degree

	// relative ke origin point
	WorldXcm Centimeter
	// relative ke origin point
	WorldYcm Centimeter
	// relative ke origin point
	WorldRcm Centimeter
	// relative ke origin point
	WorldROT Degree

	RobXcm Centimeter // cm
	RobYcm Centimeter // cm
	RobRcm Centimeter // cm
	RobROT Degree     // cm

	// Piksel relative ke omni center
	TopXpx Centimeter
	// Piksel relative ke omni center
	TopYpx Centimeter
	// Radius dalam px (euclidean dist dari camX dan camY in respect dari midpoint omni)
	TopRpx Centimeter
	// Rotasi relative ke midpoint omni, relative ke arah depan robot
	TopROT Degree
}
