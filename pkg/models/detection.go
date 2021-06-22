package models

import (
	"image"
	"math"

	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type DetectionObject struct {
	Bbox     image.Rectangle
	Midpoint image.Point
}

func EucDistance(one image.Point, other image.Point) float64 {
	xDist := math.Pow(float64(one.X-other.X), 2)
	yDist := math.Pow(float64(one.Y-other.Y), 2)

	return math.Sqrt(xDist + yDist)
}

func (d DetectionObject) AsTransform(conf *configuration.FreezeConfig) Transform {
	origin := image.Point{
		conf.Camera.MidpointRad,
		conf.Camera.MidpointRad,
	}

	// Top
	xDist := float64(d.Midpoint.X - origin.X)
	yDist := float64(origin.Y - d.Midpoint.Y)

	return Transform{
		TopXpx: Centimeter(xDist),
		TopYpx: Centimeter(yDist),
	}
}
