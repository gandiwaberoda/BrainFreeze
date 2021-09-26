package models

import (
	"math"
	"strings"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type Centimeter float64
type Degree float64
type Miliseconds int

type Radian float64

func (r Radian) AsDegree() Degree {
	return Degree(r * (180.0 / math.Pi))
}

func (r Degree) AsHalfCircle() Degree {
	y := int(r) % 360
	if y >= 0 && y <= 180 {
		return Degree(y)
	} else if y > 180 && y < 360 {
		return Degree(-180 + (y - 180))
	} else if y <= 0 && y >= -180 {
		return Degree(y)
	} else if y < -180 {
		return Degree(180 - ((-1 * y) - 180))
	}

	return Degree(r)
}
func (r Degree) ShiftRight() Degree {
	return Degree(float64(r) + 1000)
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

func (t *Transform) InjectWorldTransfromFromRobotTransform(myTransform Transform) {
	rad := float64(Degree((myTransform.WorldROT + t.RobROT - 90) * -1).AsRadian())
	x := Centimeter(math.Cos(rad) * float64(t.RobRcm))
	y := Centimeter(math.Sin(rad) * float64(t.RobRcm))

	t.WorldXcm = myTransform.WorldXcm + x
	t.WorldYcm = myTransform.WorldYcm + y
	// t.WorldRcm = t.RobRcm
}

func (t *Transform) InjectWorldTransfromFromEncTransform(conf *configuration.FreezeConfig) {
	startRot := conf.Robot.StartRot
	t.WorldROT = (t.EncROT + Degree(startRot)).AsHalfCircle() // Biar rentangnya dari -179 sampai 180

	startPos := strings.ToUpper(conf.Robot.StartPos)

	offsetX := float64(0)
	offsetY := float64(0)
	switch startPos {
	case "A":
		offsetX = 0
		offsetY = 0
	case "B":
		offsetX = 0
		offsetY = 450 / 2
	case "C":
		offsetX = 0
		offsetY = 450
	case "D":
		offsetX = 600 / 2
		offsetY = 450
	case "E":
		offsetX = 600
		offsetY = 450
	case "F":
		offsetX = 600
		offsetY = 450 / 2
	case "G":
		offsetX = 600
		offsetY = 0
	case "H":
		offsetX = 600 / 2
		offsetY = 0
	default:
		panic("Invalid startpos")
	}

	t.WorldXcm = t.EncXcm + Centimeter(offsetX)
	t.WorldYcm = t.EncYcm + Centimeter(offsetY)
}
