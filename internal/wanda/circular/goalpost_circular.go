package circular

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type GoalpostCircular struct {
	Threshold uint8
	Radius    float64
	conf      *configuration.FreezeConfig
	upperRed  gocv.Scalar
	lowerRed  gocv.Scalar
	erodeMat  gocv.Mat
}

func NewGoalpostCircular(conf *configuration.FreezeConfig) *GoalpostCircular {
	return &GoalpostCircular{
		Threshold: 253,
		Radius:    300,
		conf:      conf,
		upperRed:  gocv.NewScalar(179, 255, 255, 1),
		lowerRed:  gocv.NewScalar(165, 56, 167, 0),
		erodeMat:  gocv.Ones(9, 9, gocv.MatTypeCV8UC1),
	}
}

func (n *GoalpostCircular) Detect(hsvFrame *gocv.Mat, grayFrame *gocv.Mat) (result []models.Transform) {
	res := []models.Transform{}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERRORRR Circular goalpost")
			result = res
		}
	}()

	cBlue := color.RGBA{0, 255, 255, 0}

	hsvRed := gocv.NewMat()
	defer hsvRed.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerRed, n.upperRed, &hsvRed)
	gocv.Erode(hsvRed, &hsvRed, n.erodeMat)

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	pointVecs := gocv.FindContoursWithParams(hsvRed, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
	defer pointVecs.Close()

	detecteds := []models.DetectionObject{}
	for i := 0; i < pointVecs.Size(); i++ {
		it := pointVecs.At(i)
		area := gocv.ContourArea(it)

		if area < n.conf.Wanda.MinimumHsvArea {
			// Skip kalau ukurannya kekecilan
			continue
		}

		rect := gocv.BoundingRect(it)
		// gocv.Rectangle(&frame, rect, cRed, 2)

		d := models.NewDetectionObject(rect)
		detecteds = append(detecteds, d)
	}

	lastOneWasWhite := true
	sudutPutih := make([]float64, 0)
	for i := -179.0; i <= 180.0; i++ {
		x := int(n.Radius*math.Cos(i*math.Pi/180.0)) + n.conf.Camera.PostWidth/2
		y := int(n.Radius*math.Sin(i*math.Pi/180.0)) + n.conf.Camera.PostHeight/2

		px := grayFrame.GetUCharAt(y, x)
		if px > uint8(n.conf.Wanda.WhiteOnGrayVal) {
			if !lastOneWasWhite {
				gocv.Circle(hsvFrame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
				sudutPutih = append(sudutPutih, i)
				lastOneWasWhite = true
			}
		} else {
			lastOneWasWhite = false
		}
	}

	for i := 0; i < len(sudutPutih)-1; i++ {
		smaller := sudutPutih[i]
		bigger := sudutPutih[i+1]

		offsetsmaller := smaller + 180
		offsetsbigger := bigger + 180

		if offsetsbigger-offsetsmaller < 15 {
			continue
		}
		if offsetsbigger-offsetsmaller > 100 {
			continue
		}

		for _, v := range detecteds {
			sudut_merah := v.AsTransform(n.conf)
			if int(sudut_merah.TopRpx) > (n.conf.Camera.PostHeight - 2) {
				continue
			}

			if sudut_merah.RobROT < models.Degree(bigger) && sudut_merah.RobROT > models.Degree(smaller) {
				gocv.Line(hsvFrame, image.Point{320, 320}, v.Midpoint, cBlue, 2)
				res = append(res, sudut_merah)
			}
		}

	}

	return res
}
