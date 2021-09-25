// Terlalu berat

package radial

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopRadialDummy struct {
	MinLength   uint8
	RadialCount uint8
	conf        *configuration.FreezeConfig
}

func NewRadialDummy(conf *configuration.FreezeConfig) *TopRadialDummy {
	return &TopRadialDummy{
		MinLength:   10,
		RadialCount: 10,
		conf:        conf,
	}
}

func (n *TopRadialDummy) RadialScanLine(topHsv *gocv.Mat, deg float64) (result []models.StraightDetection) {
	res := make([]models.StraightDetection, 0)
	// gocv.Line(forHsvFrame, image.Point{n.conf.Camera.ForMidX, 0}, image.Point{n.conf.Camera.ForMidX, n.conf.Camera.ForPostHeight}, color.RGBA{255, 255, 255, 1}, 1)

	// v := pxop.GetVecbAt(*forHsvFrame, n.conf.Camera.ForPostHeight/2, n.conf.Camera.ForMidX)// gocv.GaussianBlur(*forHsvFrame, forHsvFrame, image.Point{21, 21}, 21, 21, gocv.BorderDefault)

	// pxcounter := 0            // Untuk menghitung berapa px dah kedetek
	// lastDetectedColorId := -1 // Untuk tahu apakah px berubah dari sebelumnya
	// lastStartingColorPx := 0  // Untuk tahu kapan warna yang sekarang berakhir dimulai

	xyTranslate := func(radius float64) (int, int) {
		_x := float64(n.conf.Camera.PostWidth/2) + (radius * math.Cos(deg))
		_y := float64(n.conf.Camera.PostHeight/2) + (radius * math.Sin(deg))
		// fmt.Println("xxxx: ", _x, "   ", _y)
		return int(_x), int(_y)
	}

	// finishSegment := func(curY int) {
	// 	if lastDetectedColorId != -1 && pxcounter > n.conf.Wanda.StraightMinLength {

	// 		res = append(res, models.StraightDetection{
	// 			Length:        pxcounter,
	// 			LowerY:        lastStartingColorPx,
	// 			UpperY:        curY,
	// 			DetectedColor: bfconst.ColorUsed[lastDetectedColorId],
	// 		})

	// 		// Bukan warna pertama
	// 		_col := bfconst.ColorUsed[lastDetectedColorId].Visualize
	// 		gocv.Line(forPostFrame, image.Pt(n.conf.Camera.ForMidX, lastStartingColorPx), image.Pt(n.conf.Camera.ForMidX, curY), _col, 2)
	// 		lastDetectedColorId = -1
	// 		pxcounter = 0

	// 	}
	// }

	// for r := 0; r ; y++ {
	// for r := 0; r < n.conf.Camera.PostWidth/2; r++ {
	__x, __y := xyTranslate(50)
	gocv.Circle(topHsv, image.Point{__x, __y}, 10, color.RGBA{100, 255, 255, 1}, -1)
	__x, __y = xyTranslate(100)
	gocv.Circle(topHsv, image.Point{__x, __y}, 10, color.RGBA{100, 255, 255, 1}, -1)
	__x, __y = xyTranslate(150)
	gocv.Circle(topHsv, image.Point{__x, __y}, 10, color.RGBA{100, 255, 255, 1}, -1)

	// v := pxop.GetVecbAt(*forHsvFrame, y, n.conf.Camera.ForMidX)

	// for i, accCol := range bfconst.ColorUsed {
	// 	if pxop.IsVecbInBetween(v, accCol.Upper, accCol.Lower) {

	// 		if lastDetectedColorId == i {
	// 			// Warna berlanjut
	// 			pxcounter++
	// 		} else {
	// 			// Tutup warna terakhir
	// 			finishSegment(y)

	// 			// Warna baru dimulai
	// 			lastDetectedColorId = i
	// 			pxcounter = 0
	// 			lastStartingColorPx = y
	// 		}
	// 	} else {
	// 		// Warna berhenti
	// 		finishSegment(y)

	// 		// Warna baru dimulai
	// 		// lastStartingColorPx = -1
	// 	}
	// }
	// }

	return res
}

func (n *TopRadialDummy) Detect(topHsvFrame *gocv.Mat) (result []models.StraightDetection) {
	res := make([]models.StraightDetection, 0)

	for i := -180; i < 180; i += 30 {
		rad := float64(models.Degree(i).AsRadian())
		n.RadialScanLine(topHsvFrame, rad)
	}

	return res
}
