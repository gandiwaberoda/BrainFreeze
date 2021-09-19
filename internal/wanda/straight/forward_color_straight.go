package straight

import (
	"image"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/bfconst"
	"harianugrah.com/brainfreeze/internal/wanda/pxop"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type ForwardColorStraight struct {
	MinLength uint8
	Radius    float64
	conf      *configuration.FreezeConfig
}

func NewForwardColorStraight(conf *configuration.FreezeConfig) *ForwardColorStraight {
	return &ForwardColorStraight{
		MinLength: 10,
		Radius:    150,
		conf:      conf,
	}
}

func (n *ForwardColorStraight) Detect(forHsvFrame *gocv.Mat, forPostFrame *gocv.Mat) (result []models.StraightDetection) {

	res := make([]models.StraightDetection, 0)
	// gocv.Line(forHsvFrame, image.Point{n.conf.Camera.ForMidX, 0}, image.Point{n.conf.Camera.ForMidX, n.conf.Camera.ForPostHeight}, color.RGBA{255, 255, 255, 1}, 1)

	// v := pxop.GetVecbAt(*forHsvFrame, n.conf.Camera.ForPostHeight/2, n.conf.Camera.ForMidX)// gocv.GaussianBlur(*forHsvFrame, forHsvFrame, image.Point{21, 21}, 21, 21, gocv.BorderDefault)
	gocv.Blur(*forHsvFrame, forHsvFrame, image.Point{3, 3})

	pxcounter := 0            // Untuk menghitung berapa px dah kedetek
	lastDetectedColorId := -1 // Untuk tahu apakah px berubah dari sebelumnya
	lastStartingColorPx := 0  // Untuk tahu kapan warna yang sekarang berakhir dimulai

	finishSegment := func(curY int) {
		if lastDetectedColorId != -1 && pxcounter > n.conf.Wanda.StraightMinLength {

			res = append(res, models.StraightDetection{
				Length:        pxcounter,
				LowerY:        lastStartingColorPx,
				UpperY:        curY,
				DetectedColor: bfconst.ColorUsed[lastDetectedColorId],
			})

			// Bukan warna pertama
			_col := bfconst.ColorUsed[lastDetectedColorId].Visualize
			gocv.Line(forPostFrame, image.Pt(n.conf.Camera.ForMidX, lastStartingColorPx), image.Pt(n.conf.Camera.ForMidX, curY), _col, 2)
			lastDetectedColorId = -1
			pxcounter = 0

		}
	}

	for y := n.conf.Camera.ForPostHeight; y >= 0; y-- {
		v := pxop.GetVecbAt(*forHsvFrame, y, n.conf.Camera.ForMidX)

		for i, accCol := range bfconst.ColorUsed {
			if pxop.IsVecbInBetween(v, accCol.Upper, accCol.Lower) {

				if lastDetectedColorId == i {
					// Warna berlanjut
					pxcounter++
				} else {
					// Tutup warna terakhir
					finishSegment(y)

					// Warna baru dimulai
					lastDetectedColorId = i
					pxcounter = 0
					lastStartingColorPx = y
				}
			} else {
				// Warna berhenti
				finishSegment(y)

				// Warna baru dimulai
				// lastStartingColorPx = -1
			}
		}
	}

	return res
}
