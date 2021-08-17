package bfconst

import "gocv.io/x/gocv"

var (
	DummyUpper gocv.Scalar = gocv.NewScalar(179, 255, 255, 1)
	DummyLower gocv.Scalar = gocv.NewScalar(169, 88, 181, 0)

	ForwardBallUpper gocv.Scalar = gocv.NewScalar(36, 255, 255, 1)
	ForwardBallLower gocv.Scalar = gocv.NewScalar(8, 70, 119, 0)

	MagentaUpper gocv.Scalar = gocv.NewScalar(165, 255, 255, 1)
	MagentaLower gocv.Scalar = gocv.NewScalar(144, 44, 182, 0)

	ForwardMagentaUpper gocv.Scalar = gocv.NewScalar(167, 255, 255, 1)
	ForwardMagentaLower gocv.Scalar = gocv.NewScalar(144, 44, 182, 0)
)
