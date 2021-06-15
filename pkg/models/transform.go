package models

type Centimeter float64
type Degree float64

type Transform struct {
	// cm
	encXcm Centimeter
	// cm
	encYcm Centimeter
	// derajat
	encROT Degree

	// relative ke origin point
	worldXcm Centimeter
	// relative ke origin point
	worldYcm Centimeter
	// relative ke origin point
	worldRcm Centimeter
	// relative ke origin point
	worldROT Degree

	robXcm Centimeter // cm
	robYcm Centimeter // cm
	robRcm Centimeter // cm
	robROT Degree     // cm

	// Piksel relative ke omni center
	topXpx Centimeter
	// Piksel relative ke omni center
	topYpx Centimeter
	// Radius dalam px (euclidean dist dari camX dan camY in respect dari midpoint omni)
	topRpx Centimeter
	// Rotasi relative ke midpoint omni, relative ke arah depan robot
	topROT Degree
}
