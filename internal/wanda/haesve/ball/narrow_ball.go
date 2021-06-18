package ball

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type NarrowHaesveBall struct {
	conf     *configuration.FreezeConfig
	upperHsv gocv.Scalar
	lowerHsv gocv.Scalar
}

func NewNarrowHaesveBall(conf *configuration.FreezeConfig) *NarrowHaesveBall {
	return &NarrowHaesveBall{
		conf: conf,
	}
}

func getRectMidpoint(rect image.Rectangle) image.Point {
	x := (rect.Max.X + rect.Min.X) / 2
	y := (rect.Max.Y + rect.Min.Y) / 2

	return image.Pt(x, y)
}

func filter(src gocv.Mat, dst *gocv.Mat) {
	upper := gocv.NewScalar(33, 255, 255, 1)
	lower := gocv.NewScalar(8, 120, 112, 0)

	gocv.InRangeWithScalar(src, lower, upper, dst)
}

// Input adalah Mat yang sudah dalam format hsv
func (n *NarrowHaesveBall) Detect(hsvFrame gocv.Mat) []models.DetectionObject {
	w := hsvFrame.Cols()
	h := hsvFrame.Rows()

	filtered := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC1)

	filter(hsvFrame, &filtered)

	erodeMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	gocv.Erode(filtered, &filtered, erodeMat)

	dilateMat := gocv.Ones(21, 21, gocv.MatTypeCV8UC1)
	gocv.Dilate(filtered, &filtered, dilateMat)

	c := color.RGBA{255, 0, 0, 0}

	hierarchyMat := gocv.NewMat()
	pointVecs := gocv.FindContoursWithParams(filtered, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)

	detecteds := []models.DetectionObject{}
	for i := 0; i < pointVecs.Size(); i++ {
		it := pointVecs.At(i)
		area := gocv.ContourArea(it)
		fmt.Println("x", area)
		if area < n.conf.Wanda.MinimumHsvArea {
			// Skip kalau ukurannya kekecilan
			continue
		}

		rect := gocv.BoundingRect(it)
		gocv.Rectangle(&hsvFrame, rect, c, 2)
		fmt.Println(rect)

		d := models.DetectionObject{
			Bbox:     rect,
			Midpoint: getRectMidpoint(rect),
		}
		detecteds = append(detecteds, d)
	}

	return detecteds
}
