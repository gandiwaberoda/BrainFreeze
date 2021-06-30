package cyan

import (
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type NarrowHaesveCyan struct {
	conf     *configuration.FreezeConfig
	upperHsv gocv.Scalar
	lowerHsv gocv.Scalar
}

func NewNarrowHaesveCyan(conf *configuration.FreezeConfig) *NarrowHaesveCyan {
	upper := gocv.NewScalar(113, 255, 255, 1)
	lower := gocv.NewScalar(95, 88, 65, 0)

	return &NarrowHaesveCyan{
		conf:     conf,
		upperHsv: upper,
		lowerHsv: lower,
	}
}

// Input adalah Mat yang sudah dalam format hsv
func (n *NarrowHaesveCyan) Detect(hsvFrame *gocv.Mat) (bool, []models.DetectionObject) {
	detecteds := []models.DetectionObject{}

	if hsvFrame.Empty() {
		return false, detecteds
	}

	// w := hsvFrame.Cols()
	// h := hsvFrame.Rows()

	// filtered := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC1)
	filtered := gocv.NewMat()
	defer filtered.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerHsv, n.upperHsv, &filtered)

	erodeMat := gocv.Ones(4, 4, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()
	gocv.Erode(filtered, &filtered, erodeMat)

	dilateMat := gocv.Ones(17, 17, gocv.MatTypeCV8UC1)
	defer dilateMat.Close()
	gocv.Dilate(filtered, &filtered, dilateMat)

	c := color.RGBA{0, 128, 128, 0}

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

		if area < n.conf.Wanda.MinimumHsvArea {
			// Skip kalau ukurannya kekecilan
			continue
		}
		if area > n.conf.Wanda.MaximumHsvArea {
			continue
		}

		rect := gocv.BoundingRect(it)
		gocv.Rectangle(hsvFrame, rect, c, 2)
		gocv.PutText(hsvFrame, "Cyan", rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)

		d := models.NewDetectionObject(rect)
		detecteds = append(detecteds, d)
	}

	return true, detecteds
}
