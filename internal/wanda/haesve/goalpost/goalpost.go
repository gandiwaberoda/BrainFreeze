package goalpost

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type HaesveGoalpost struct {
	conf     *configuration.FreezeConfig
	upperHsv gocv.Scalar
	lowerHsv gocv.Scalar
}

func NewHaesveGoalpost(conf *configuration.FreezeConfig) *HaesveGoalpost {
	upper := gocv.NewScalar(166, 51, 255, 1)
	lower := gocv.NewScalar(0, 0, 255, 0)

	return &HaesveGoalpost{
		conf:     conf,
		upperHsv: upper,
		lowerHsv: lower,
	}
}

// Input adalah Mat yang sudah dalam format hsv
func (n *HaesveGoalpost) Detect(hsvFrame *gocv.Mat) (found bool, result []models.DetectionObject) {
	detecteds := []models.DetectionObject{}

	defer func() {
		if r := recover(); r != nil {
			found = false
			result = detecteds
			fmt.Println("recovered from ", r)
			return
		}
	}()

	if hsvFrame.Empty() {
		return false, detecteds
	}

	w := hsvFrame.Cols()
	h := hsvFrame.Rows()

	filtered := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC1)
	defer filtered.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerHsv, n.upperHsv, &filtered)

	// dilateMat := gocv.Ones(9, 9, gocv.MatTypeCV8UC1)
	// defer dilateMat.Close()
	// gocv.Dilate(filtered, &filtered, dilateMat)

	// erodeMat := gocv.Ones(35, 35, gocv.MatTypeCV8UC1)
	// defer erodeMat.Close()
	// gocv.Erode(filtered, &filtered, erodeMat)

	c := color.RGBA{128, 0, 128, 0}

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	pointVecs := gocv.FindContoursWithParams(filtered, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
	defer pointVecs.Close()

	if pointVecs.Size() == 0 {
		return false, detecteds
	}

	for i := 0; i < pointVecs.Size(); i++ {
		it := pointVecs.At(i)
		area := gocv.ContourArea(it)

		// if area < n.conf.Wanda.MinimumHsvArea {
		// 	// Skip kalau ukurannya kekecilan
		// 	continue
		// }
		// if area > n.conf.Wanda.MaximumHsvArea {
		// 	continue
		// }

		rect := gocv.BoundingRect(it)
		gocv.Rectangle(hsvFrame, rect, c, 2)
		gocv.PutText(hsvFrame, "Goalpost", rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)

		d := models.NewDetectionObject(rect)
		d.ContourArea = area
		d.BboxArea = float64((rect.Max.X - rect.Min.X) * (rect.Max.Y - rect.Min.Y))
		detecteds = append(detecteds, d)
	}

	return true, detecteds
}
