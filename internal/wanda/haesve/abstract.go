package haesve

import (
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
)

type HaesveInterface interface {
	Detect(gocv.Mat) []models.DetectionObject
}
