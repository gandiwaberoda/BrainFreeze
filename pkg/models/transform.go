package models

import "math"

type Centimeter float64
type Degree float64
type Miliseconds int

type Radian float64

func (r Radian) AsDegree() Degree {
	return Degree(r * (180.0 / math.Pi))
}

func (d Degree) AsRadian() Radian {
	return Radian(d * (math.Pi / 180.0))
}

func (d *Degree) Rotate(am Degree) {
	amount := math.Mod(float64(am), 360.0)
	v := float64(*d) + amount

	if v > 180 {
		v = -1 * (360 - v)
	}

	if v < -180 {
		v = -1 * (-360 - v)
	}

	*d = Degree(v)

	if v < -360 || v > 360 {
		d.Rotate(0)
	}
}

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
