package models

import (
	"image"
	"math"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type DetectionObject struct {
	Bbox        image.Rectangle
	Midpoint    image.Point
	OuterRad    int
	CloserPoint image.Point
}

func NewDetectionObject(bbox image.Rectangle) DetectionObject {
	d := DetectionObject{Bbox: bbox}

	xMid := (d.Bbox.Max.X + d.Bbox.Min.X) / 2
	yMid := (d.Bbox.Max.Y + d.Bbox.Min.Y) / 2
	d.Midpoint = image.Pt(xMid, yMid)

	outerCicle := d.Bbox.Max.X - d.Midpoint.X
	if vert := d.Bbox.Max.X - d.Midpoint.X; vert > outerCicle {
		outerCicle = vert
	}
	d.OuterRad = outerCicle

	// TODO: Pake titik pertemuan lingkaran dan garis terdekat
	d.CloserPoint = d.Midpoint

	return d
}

func EucDistance(xDist, yDist float64) float64 {
	// xDist := math.Pow(float64(one.X-other.X), 2)
	// yDist := math.Pow(float64(one.Y-other.Y), 2)

	xD2 := math.Pow(xDist, 2)
	yD2 := math.Pow(yDist, 2)

	return math.Sqrt(xD2 + yD2)
}

func (d DetectionObject) AsTransform(conf *configuration.FreezeConfig) Transform {
	// Top
	xTopDist := float64(d.Midpoint.X - conf.Camera.Midpoint.X)
	yTopDist := float64(conf.Camera.Midpoint.Y - d.Midpoint.Y)
	rTopDist := EucDistance(xTopDist, yTopDist)

	rotTopDegree := Radian(math.Atan2(yTopDist, xTopDist)).AsDegree()
	// Balik
	rotTopDegree = rotTopDegree * -1

	return Transform{
		TopXpx: Centimeter(xTopDist),
		TopYpx: Centimeter(yTopDist),
		TopRpx: Centimeter(rTopDist),
		TopROT: rotTopDegree,
	}
}
