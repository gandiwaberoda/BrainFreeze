package bfconst

import "gocv.io/x/gocv"

var (
	DummyUpper gocv.Scalar = gocv.NewScalar(179, 255, 255, 1)
	DummyLower gocv.Scalar = gocv.NewScalar(169, 88, 181, 0)

	ForwardBallUpper gocv.Scalar = gocv.NewScalar(36, 255, 255, 1)
	ForwardBallLower gocv.Scalar = gocv.NewScalar(8, 70, 119, 0)

	ForwardMagentaUpper gocv.Scalar = gocv.NewScalar(167, 255, 255, 1)
	ForwardMagentaLower gocv.Scalar = gocv.NewScalar(149, 34, 149, 0)
)
