package circular

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/bfconst"
	"harianugrah.com/brainfreeze/internal/wanda/pxop"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type DummyAvoidCircular struct {
	FOVMin float64
	FOVMax float64
	Radius float64
	conf   *configuration.FreezeConfig
}

func NewDummyAvoidCircular(conf *configuration.FreezeConfig) *DummyAvoidCircular {
	return &DummyAvoidCircular{
		FOVMin: -179,
		FOVMax: 180,
		// Radius: 150,
		Radius: float64(conf.Wanda.RadiusDummyCircular),
		conf:   conf,
	}
}

func (n *DummyAvoidCircular) Detect(tophsvFrame *gocv.Mat) (result []float64) {
	res := make([]float64, 0)

	for i := n.FOVMin; i <= n.FOVMax; i++ {
		x := int(n.Radius*math.Cos(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostWidth/2
		y := int(n.Radius*math.Sin(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostHeight/2

		hsv := pxop.GetVecbAt(*tophsvFrame, y, x)
		gocv.Circle(tophsvFrame, image.Point{X: x, Y: y}, 5, color.RGBA{0, 0, 0, 0}, -1)
		if pxop.IsVecbInBetween(hsv, bfconst.Dummy.Upper, bfconst.Dummy.Lower) {
			gocv.Circle(tophsvFrame, image.Point{X: x, Y: y}, 10, color.RGBA{0, 0, 255, 1}, -1)
			res = append(res, i)
		}
	}

	return res
}
