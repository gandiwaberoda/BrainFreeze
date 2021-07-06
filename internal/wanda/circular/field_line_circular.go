package circular

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type FieldLineCircular struct {
	FOVMin    float64
	FOVMax    float64
	Threshold uint8
	Radius    float64
	conf      *configuration.FreezeConfig
}

func NewFieldLineCircular(conf *configuration.FreezeConfig) *FieldLineCircular {
	return &FieldLineCircular{
		FOVMin:    -120,
		FOVMax:    120,
		Threshold: 253,
		Radius:    150,
		conf:      conf,
	}
}

func (n *FieldLineCircular) Detect(grayFrame *gocv.Mat) (result []float64) {
	res := make([]float64, 0)

	lastOneWasWhite := false
	for i := n.FOVMin; i <= n.FOVMax; i++ {
		x := int(n.Radius*math.Cos(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostWidth/2
		y := int(n.Radius*math.Sin(float64(i)*math.Pi/180.0)) + n.conf.Camera.PostHeight/2

		px := grayFrame.GetUCharAt(y, x)

		// gocv.Circle(grayFrame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
		if px > n.Threshold {
			if !lastOneWasWhite {
				gocv.Circle(grayFrame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
				res = append(res, i)
				lastOneWasWhite = true
			}
		} else {
			lastOneWasWhite = false
		}
	}

	// fmt.Println(len(res))
	return res
}
