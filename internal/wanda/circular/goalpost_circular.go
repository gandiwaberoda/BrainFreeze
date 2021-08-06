package circular

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/bfconst"
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
		upperRed:  bfconst.DummyUpper,
		lowerRed:  bfconst.DummyLower,
		erodeMat:  gocv.Ones(5, 5, gocv.MatTypeCV8UC1),
	}
}

func (n *GoalpostCircular) Detect(hsvFrame *gocv.Mat, grayFrame *gocv.Mat) (result []models.DetectionObject) {
	res := []models.DetectionObject{}

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("ERRORRR Circular goalpost")
			result = res
		}
	}()

	hsvRed := gocv.NewMat()
	defer hsvRed.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerRed, n.upperRed, &hsvRed)
	// gocv.Erode(hsvRed, &hsvRed, n.erodeMat)

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	pointVecs := gocv.FindContoursWithParams(hsvRed, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
	defer pointVecs.Close()

	detecteds := []models.DetectionObject{}
	for i := 0; i < pointVecs.Size(); i++ {
		it := pointVecs.At(i)
		area := gocv.ContourArea(it)

		if area < 300 {
			// Skip kalau ukurannya kekecilan
			continue
		}

		rect := gocv.BoundingRect(it)
		d := models.NewDetectionObject(rect)
		dist := models.EucDistance(float64(d.Midpoint.X)-320, float64(d.Midpoint.Y)-320)
		fmt.Println("dist", dist)
		if dist < 70 {
			continue
		}
		gocv.Rectangle(hsvFrame, rect, color.RGBA{0, 0, 0, 1}, 2)

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
				gocv.Circle(hsvFrame, image.Point{X: x, Y: y}, 5, color.RGBA{0, 0, 0, 0}, -1)
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

				v.BboxArea = float64((v.Bbox.Max.X - v.Bbox.Min.X) * (v.Bbox.Max.Y - v.Bbox.Min.Y))
				res = append(res, v)
			}
		}

	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].BboxArea > res[j].BboxArea
	})

	return res
}
