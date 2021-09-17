package straight

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/bfconst"
	"harianugrah.com/brainfreeze/internal/wanda/pxop"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type ForwardColorStraight struct {
	FOVMin    float64
	FOVMax    float64
	Threshold uint8
	Radius    float64
	conf      *configuration.FreezeConfig
}

func NewForwardColorStraight(conf *configuration.FreezeConfig) *ForwardColorStraight {
	return &ForwardColorStraight{
		FOVMin:    float64(conf.Wanda.LfFovMin),
		FOVMax:    float64(conf.Wanda.LfFovMax),
		Threshold: 253,
		Radius:    150,
		conf:      conf,
	}
}

func (n *ForwardColorStraight) Detect(forHsvFrame *gocv.Mat, forPostFrame *gocv.Mat) (result []float64) {

	res := make([]float64, 0)
	// gocv.Line(forHsvFrame, image.Point{n.conf.Camera.ForMidX, 0}, image.Point{n.conf.Camera.ForMidX, n.conf.Camera.ForPostHeight}, color.RGBA{255, 255, 255, 1}, 1)

	// v := pxop.GetVecbAt(*forHsvFrame, n.conf.Camera.ForPostHeight/2, n.conf.Camera.ForMidX)
	gocv.Circle(forPostFrame, image.Point{n.conf.Camera.ForMidX, n.conf.Camera.ForPostHeight / 2}, 7, color.RGBA{255, 255, 255, 1}, 1)

	for y := 0; y < n.conf.Camera.ForPostHeight; y++ {
		v := pxop.GetVecbAt(*forHsvFrame, y, n.conf.Camera.ForMidX)

		// Dummy
		if pxop.IsVecbInBetween(v, pxop.VecbFrom4(bfconst.DummyUpper), pxop.VecbFrom4(bfconst.DummyLower)) {
			gocv.Circle(forPostFrame, image.Point{n.conf.Camera.ForMidX, y}, 2, color.RGBA{255, 0, 0, 1}, 1)
		}

		// Bola
		if pxop.IsVecbInBetween(v, pxop.VecbFrom4(bfconst.ForwardBallUpper), pxop.VecbFrom4(bfconst.ForwardBallLower)) {
			gocv.Circle(forPostFrame, image.Point{n.conf.Camera.ForMidX, y}, 4, color.RGBA{250, 190, 88, 1}, 1)
		}

		// Magenta
		if pxop.IsVecbInBetween(v, pxop.VecbFrom4(bfconst.MagentaUpper), pxop.VecbFrom4(bfconst.MagentaLower)) {
			gocv.Circle(forPostFrame, image.Point{n.conf.Camera.ForMidX, y}, 6, color.RGBA{255, 0, 255, 1}, 1)
		}

		// Cyan
		if pxop.IsVecbInBetween(v, pxop.VecbFrom4(bfconst.CyanUpper), pxop.VecbFrom4(bfconst.CyanLower)) {
			gocv.Circle(forPostFrame, image.Point{n.conf.Camera.ForMidX, y}, 7, color.RGBA{0, 0, 255, 1}, 1)
		}

	}

	// lastOneWasWhite := false
	// for i := n.FOVMin; i <= n.FOVMax; i++ {
	// 	x := int(n.Radius*math.Cos(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostWidth/2
	// 	y := int(n.Radius*math.Sin(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostHeight/2

	// 	px := grayFrame.GetUCharAt(y, x)

	// 	// gocv.Circle(grayFrame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
	// 	if px > n.Threshold {
	// 		if !lastOneWasWhite {
	// 			gocv.Circle(grayFrame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
	// 			res = append(res, i)
	// 			lastOneWasWhite = true
	// 		}
	// 	} else {
	// 		lastOneWasWhite = false
	// 	}
	// }

	// fmt.Println(len(res))
	return res
}
