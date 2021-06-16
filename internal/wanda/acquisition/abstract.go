package acquisition

import "gocv.io/x/gocv"

type CameraAcquisitionInterface interface {
	Start() (bool, error)
	Stop() (bool, error)
	Read() gocv.Mat
}
