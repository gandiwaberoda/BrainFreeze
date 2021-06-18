package models

import "image"

type DetectionObject struct {
	Bbox     image.Rectangle
	Midpoint image.Point
}

func (d DetectionObject) AsTransform() Transform {
	return Transform{
		EncXcm: 10,
		EncYcm: 10,
		EncROT: 30,
	}
}
