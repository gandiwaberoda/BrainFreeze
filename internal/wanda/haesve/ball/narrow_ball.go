package ball

import (
	// "fmt"
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
	lower := gocv.NewScalar(6, 120, 70, 0)

	return &NarrowHaesveBall{
		conf:     conf,
		upperHsv: upper,
		lowerHsv: lower,
	}
}

// Input adalah Mat yang sudah dalam format hsv
func (n *NarrowHaesveBall) Detect(hsvFrame *gocv.Mat) []models.DetectionObject {
	// defer hsvFrame.Close()
	// win := gocv.NewWindow("tanya")
	// win.IMShow(*hsvFrame)
	// win.WaitKey(1)

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
		gocv.Rectangle(hsvFrame, rect, c, 2)
		fmt.Println(rect)

		// d := models.DetectionObject{
		// 	Bbox:     rect,
		// 	Midpoint: getRectMidpoint(rect),
		// }
		d := models.NewDetectionObject(rect)
		detecteds = append(detecteds, d)
	}

	return detecteds
}
