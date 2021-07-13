package main

import (
	"image"
	"image/color"
	"math"
)

func IsWhite(ac color.Color) bool {
	max := uint32(65535)
	r, g, b, _ := ac.RGBA()
	return r == max && g == max && b == max
}

type LidarReading struct {
	ClosestPoint float64
}
type SensorModel struct {
	numSensor   int
	minRotation float64
	maxRotation float64
	deadZoneRad float64
	maxDistance float64
	Reading     map[float64]LidarReading
}

func NewSensorModel() SensorModel {
	reading := make(map[float64]LidarReading)
	min := -130.0
	max := 130.0
	num := 14
	sen_range := math.Abs(min) + math.Abs(max)
	sen_step := sen_range / float64(num)

	for i := min; i <= max; i += sen_step {
		reading[i] = LidarReading{}
	}

	return SensorModel{
		numSensor:   num,
		minRotation: min,
		maxRotation: max,
		maxDistance: 320,
		Reading:     reading,
		deadZoneRad: 10,
	}
}

func (rob *SensorModel) SenseFromImage(img image.Image, worldPos image.Point) {
	for k := range rob.Reading {
		newReading := LidarReading{}
		for r := rob.deadZoneRad; r <= rob.maxDistance; r++ {
			endX, endY := Polar2Cartesian(k, r)
			endX += float64(worldPos.X)
			endY += float64(worldPos.Y)

			px := img.At(int(endX), int(endY))
			if IsWhite(px) {
				newReading.ClosestPoint = r
				break
			}
		}
		rob.Reading[k] = newReading
	}
}

func Polar2Cartesian(deg, rad float64) (x, y float64) {
	radian := ((deg * -1) - 90) * math.Pi / 180
	x = rad * math.Cos(radian)
	y = rad * math.Sin(radian)
	return
}
