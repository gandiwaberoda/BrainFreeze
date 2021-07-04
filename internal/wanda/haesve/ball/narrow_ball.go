package ball

import (
	"fmt"
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
	upper := gocv.NewScalar(37, 255, 255, 1)
	lower := gocv.NewScalar(6, 150, 70, 0)

	return &NarrowHaesveBall{
		conf:     conf,
		upperHsv: upper,
		lowerHsv: lower,
	}
}

// Input adalah Mat yang sudah dalam format hsv
func (n *NarrowHaesveBall) Detect(hsvFrame *gocv.Mat) (found bool, result []models.DetectionObject) {
	detecteds := []models.DetectionObject{}

	defer func() {
		if r := recover(); r != nil {
			found = false
			result = detecteds
			fmt.Println("recovered from ", r)
			return
		}
	}()

	// defer hsvFrame.Close()
	// win := gocv.NewWindow("tanya")
	// win.IMShow(*hsvFrame)
	// win.WaitKey(1)

	if hsvFrame.Empty() {
		return false, detecteds
	}

	w := hsvFrame.Cols()
	h := hsvFrame.Rows()

	filtered := gocv.NewMatWithSize(h, w, gocv.MatTypeCV8UC1)
	defer filtered.Close()

	gocv.InRangeWithScalar(*hsvFrame, n.lowerHsv, n.upperHsv, &filtered)

	erodeMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()
	gocv.Erode(filtered, &filtered, erodeMat)

	dilateMat := gocv.Ones(17, 17, gocv.MatTypeCV8UC1)
	defer dilateMat.Close()
	gocv.Dilate(filtered, &filtered, dilateMat)

	c := color.RGBA{255, 0, 0, 0}

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	pointVecs := gocv.FindContoursWithParams(filtered, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
	defer pointVecs.Close()

	if pointVecs.Size() == 0 {
		return false, detecteds
	}

	for i := 0; i < pointVecs.Size(); i++ {
		// <-time.After(time.Millisecond * 1500)
		it := pointVecs.At(i)
		area := gocv.ContourArea(it)

		// fmt.Println("ENTAH:", i, "AREA:", area)

		if area < n.conf.Wanda.MinimumHsvArea {
			// Skip kalau ukurannya kekecilan
			continue
		}
		if area > n.conf.Wanda.MaximumHsvArea {
			continue
		}

		rect := gocv.BoundingRect(it)
		gocv.Rectangle(hsvFrame, rect, c, 2)

		d := models.NewDetectionObject(rect)
		detecteds = append(detecteds, d)
	}

	return true, detecteds
}
