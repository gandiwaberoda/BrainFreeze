package haesve

import (
	"image"

	"gocv.io/x/gocv"
)

type HaesevDetected struct {
	image.Rectangle
	midpoint image.Point
}

type HaesveInterface interface {
	Detect(gocv.Mat) []HaesevDetected
}
