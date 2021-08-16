package ball

// import (
// 	"fmt"

// 	"gocv.io/x/gocv"
// 	"harianugrah.com/brainfreeze/internal/bfconst"
// 	"harianugrah.com/brainfreeze/pkg/models"
// 	"harianugrah.com/brainfreeze/pkg/models/configuration"
// )

// type ForwardNarrowHaesveBall struct {
// 	conf     *configuration.FreezeConfig
// 	upperHsv gocv.Scalar
// 	lowerHsv gocv.Scalar
// }

// func NewForwardNarrowHaesveMagenta(conf *configuration.FreezeConfig) *ForwardNarrowHaesveBall {
// 	return &ForwardNarrowHaesveBall{
// 		conf:     conf,
// 		upperHsv: bfconst.ForwardBallUpper,
// 		lowerHsv: bfconst.ForwardBallLower,
// 	}
// }

// // Input adalah Mat yang sudah dalam format hsv
// func (n *ForwardNarrowHaesveBall) Detect(hsvFrame *gocv.Mat) (found bool, result []models.DetectionObject) {
// 	detecteds := []models.DetectionObject{}

// 	defer func() {
// 		if r := recover(); r != nil {
// 			found = false
// 			result = detecteds
// 			fmt.Println("recovered from ", r)
// 			return
// 		}
// 	}()

// 	if hsvFrame.Empty() {
// 		return false, detecteds
// 	}
// 	filtered := gocv.NewMat()
// 	defer filtered.Close()

// 	gocv.InRangeWithScalar(*hsvFrame, n.lowerHsv, n.upperHsv, &filtered)

// 	erodeMat := gocv.Ones(15, 15, gocv.MatTypeCV8UC1)
// 	defer erodeMat.Close()
// 	gocv.Erode(filtered, &filtered, erodeMat)

// 	dilateMat := gocv.Ones(7, 7, gocv.MatTypeCV8UC1)
// 	defer dilateMat.Close()
// 	gocv.Dilate(filtered, &filtered, dilateMat)

// 	hierarchyMat := gocv.NewMat()
// 	defer hierarchyMat.Close()

// 	pointVecs := gocv.FindContoursWithParams(filtered, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
// 	defer pointVecs.Close()

// 	if pointVecs.Size() == 0 {
// 		return false, detecteds
// 	}

// 	for i := 0; i < pointVecs.Size(); i++ {
// 		it := pointVecs.At(i)
// 		area := gocv.ContourArea(it)

// 		if area < n.conf.Wanda.MinimumHsvArea {
// 			// Skip kalau ukurannya kekecilan
// 			continue
// 		}

// 		rect := gocv.BoundingRect(it)
// 		// gocv.Rectangle(hsvFrame, rect, c, 2)
// 		// gocv.PutText(hsvFrame, "Magenta Forward", rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)

// 		d := models.NewDetectionObject(rect)
// 		detecteds = append(detecteds, d)
// 	}

// 	return true, detecteds
// }
