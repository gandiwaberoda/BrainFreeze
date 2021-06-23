package models

import (
	"image"
	"math"
	"sort"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type DetectionObject struct {
	Bbox        image.Rectangle
	Midpoint    image.Point
	OuterRad    int
	CloserPoint image.Point
}

func (d DetectionObject) Lerp(other DetectionObject, percentage float64) DetectionObject {
	x := (percentage * float64(other.Midpoint.X)) + ((1 - percentage) * float64(d.Midpoint.X))
	y := (percentage * float64(other.Midpoint.Y)) + ((1 - percentage) * float64(d.Midpoint.Y))

	d.Midpoint = image.Point{
		X: int(x),
		Y: int(y),
	}

	return d
}

// Untuk mencari titik yang paling masuk akal menjadi bola, jika diketahui lokasi bola sebelumnya
func (d DetectionObject) SortDetectionsObjectByDistanceToMe(other []DetectionObject) []DetectionObject {
	sort.Slice(other, func(i, j int) bool {
		distI := EucDistance(float64(d.Midpoint.X)-float64(other[i].Midpoint.X), float64(d.Midpoint.Y)-float64(other[i].Midpoint.Y))
		distJ := EucDistance(float64(d.Midpoint.X)-float64(other[j].Midpoint.X), float64(d.Midpoint.Y)-float64(other[j].Midpoint.Y))

		return distI < distJ
	})

	return other
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
	xD2 := math.Pow(xDist, 2)
	yD2 := math.Pow(yDist, 2)

	return math.Sqrt(xD2 + yD2)
}

func pxToCm(px float64) float64 {
	return px * 0.6
}

func (d DetectionObject) AsTransform(conf *configuration.FreezeConfig) Transform {
	// Top
	xTopDist := float64(d.Midpoint.X - conf.Camera.Midpoint.X)
	yTopDist := float64(conf.Camera.Midpoint.Y - d.Midpoint.Y)
	rTopDist := EucDistance(xTopDist, yTopDist)

	robROTDegree := Radian(math.Atan2(yTopDist, xTopDist)).AsDegree()
	// Balik
	robROTDegree = robROTDegree * -1
	robROTDegree.Rotate(0)

	rotTopDegree := robROTDegree
	rotTopDegree.Rotate(Degree(conf.Camera.RobFrontOffsetDeg))

	robRcm := pxToCm(rTopDist)
	robXcm := robRcm * math.Sin(float64(robROTDegree.AsRadian()))
	robYcm := robRcm * math.Cos(float64(robROTDegree.AsRadian()))

	return Transform{
		TopXpx: Centimeter(xTopDist),
		TopYpx: Centimeter(yTopDist),
		TopRpx: Centimeter(rTopDist),
		TopROT: rotTopDegree,

		RobXcm: Centimeter(pxToCm(robXcm)),
		RobYcm: Centimeter(pxToCm(robYcm)),
		RobRcm: Centimeter(robRcm),
		RobROT: robROTDegree,
	}
}
