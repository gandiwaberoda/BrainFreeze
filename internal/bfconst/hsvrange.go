package bfconst

import (
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/pxop"
	"harianugrah.com/brainfreeze/pkg/models"
)

var (
	// Dummy
	DummyUpper gocv.Scalar = gocv.NewScalar(179, 255, 255, 1)
	DummyLower gocv.Scalar = gocv.NewScalar(169, 88, 181, 0)

	// Bola
	BallUpper        gocv.Scalar = gocv.NewScalar(37, 255, 255, 1)
	BallLower        gocv.Scalar = gocv.NewScalar(0, 97, 193, 0)
	ForwardBallUpper gocv.Scalar = gocv.NewScalar(36, 255, 255, 1)
	ForwardBallLower gocv.Scalar = gocv.NewScalar(8, 70, 119, 0)

	// Magenta
	MagentaUpper        gocv.Scalar = gocv.NewScalar(171, 255, 255, 1)
	MagentaLower        gocv.Scalar = gocv.NewScalar(144, 85, 170, 0)
	ForwardMagentaUpper gocv.Scalar = gocv.NewScalar(167, 255, 255, 1)
	ForwardMagentaLower gocv.Scalar = gocv.NewScalar(144, 44, 182, 0)

	// Cyan
	CyanUpper gocv.Scalar = gocv.NewScalar(124, 220, 255, 1)
	CyanLower gocv.Scalar = gocv.NewScalar(93, 82, 128, 0)
)

var (
	Dummy   models.AcceptableColor = models.AcceptableColor{Id: 1, Name: "Dummy", Upper: pxop.VecbFrom4(DummyUpper), Lower: pxop.VecbFrom4(DummyLower), Visualize: color.RGBA{255, 0, 0, 1}}
	Ball    models.AcceptableColor = models.AcceptableColor{Id: 2, Name: "Ball", Upper: pxop.VecbFrom4(BallUpper), Lower: pxop.VecbFrom4(BallLower), Visualize: color.RGBA{250, 190, 88, 1}}
	Magenta models.AcceptableColor = models.AcceptableColor{Id: 3, Name: "Magenta", Upper: pxop.VecbFrom4(MagentaUpper), Lower: pxop.VecbFrom4(MagentaLower), Visualize: color.RGBA{255, 0, 255, 1}}
	Cyan    models.AcceptableColor = models.AcceptableColor{Id: 4, Name: "Cyan", Upper: pxop.VecbFrom4(CyanUpper), Lower: pxop.VecbFrom4(CyanLower), Visualize: color.RGBA{0, 0, 255, 1}}
)

var ColorUsed = []models.AcceptableColor{
	Dummy, Ball, Magenta, Cyan,
}
