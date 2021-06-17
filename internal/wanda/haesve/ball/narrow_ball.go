package ball

import (
	"fmt"
	"image/color"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/haesve"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type NarrowHaesveBall struct {
	conf configuration.FreezeConfig
}

func CreateNewNarrowHaesveBall(conf configuration.FreezeConfig) *NarrowHaesveBall {
	return &NarrowHaesveBall{
		conf: conf,
	}
}

func filter(src gocv.Mat, dst *gocv.Mat) {
	upper := gocv.NewScalar(33, 255, 255, 1)
	lower := gocv.NewScalar(8, 120, 112, 0)

	gocv.InRangeWithScalar(src, lower, upper, dst)
}

// Input adalah Mat yang sudah dalam format hsv
func (n *NarrowHaesveBall) Detect(hsvFrame gocv.Mat) []haesve.HaesevDetected {
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
	for i := 0; i < pointVecs.Size(); i++ {
		it := pointVecs.At(i)
		// area := gocv.ContourArea(it)
		// if area < 5 {
		// 	// Skip kalau ukurannya kekecilan
		// 	continue
		// }

		rect := gocv.BoundingRect(it)
		gocv.Rectangle(&hsvFrame, rect, c, 2)
		fmt.Println(rect)
	}

	gocv.IMWrite("./temp/"+time.Now().String()+".jpg", hsvFrame)

	return nil
}
