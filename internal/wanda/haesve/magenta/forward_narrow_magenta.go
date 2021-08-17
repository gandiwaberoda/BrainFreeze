package magenta

import (
	"fmt"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/bfconst"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type ForwardNarrowHaesveMagenta struct {
	conf     *configuration.FreezeConfig
	upperHsv gocv.Scalar
	lowerHsv gocv.Scalar
}

func NewForwardNarrowHaesveMagenta(conf *configuration.FreezeConfig) *ForwardNarrowHaesveMagenta {
	return &ForwardNarrowHaesveMagenta{
		conf:     conf,
		upperHsv: bfconst.ForwardMagentaUpper,
		lowerHsv: bfconst.ForwardMagentaLower,
	}
}

// Input adalah Mat yang sudah dalam format hsv
func (n *ForwardNarrowHaesveMagenta) Detect(hsvFrame *gocv.Mat) (found bool, result []models.DetectionObject) {
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
	filtered := gocv.NewMat()
	defer filtered.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerHsv, n.upperHsv, &filtered)

	// erodeMat := gocv.Ones(5, 5, gocv.MatTypeCV8UC1)
	// defer erodeMat.Close()
	// gocv.Erode(filtered, &filtered, erodeMat)

	nm := gocv.NewMat()
	defer nm.Close()

	// dilateMat := gocv.Ones(9, 9, gocv.MatTypeCV8UC1)
	// defer dilateMat.Close()
	// gocv.Dilate(filtered, &nm, dilateMat)

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

		rect := gocv.BoundingRect(it)

		d := models.NewDetectionObject(rect)
		d.ContourArea = area
		detecteds = append(detecteds, d)
	}

	return true, detecteds
}
