package wanda

import (
	"image/color"

	// "gocv.io/x/gocv"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type WandaVision struct {
	isRunning  bool
	conf       *configuration.FreezeConfig
	topCamera  *acquisition.TopCameraAcquisition
	ballNarrow *ball.NarrowHaesveBall
	state      *state.StateAccess
}

func NewWandaVision(conf *configuration.FreezeConfig, state *state.StateAccess) *WandaVision {
	return &WandaVision{
		conf:  conf,
		state: state,
	}
}

func worker(w *WandaVision) {
	win := gocv.NewWindow("hhh")

	warna := color.RGBA{0, 255, 0, 0}
	hsvFrame := gocv.NewMat()

	frame := gocv.NewMat()
	for {
		w.topCamera.Read(&frame)

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)

		// Ball
		narrowBallRes := w.ballNarrow.Detect(&hsvFrame)
		if len(narrowBallRes) > 0 {
			transform := narrowBallRes[0].AsTransform(w.conf)
			gocv.Rectangle(&frame, narrowBallRes[0].Bbox, warna, 3)
			w.state.UpdateBallTransform(transform)
		}

		// EGP

		// FGP

		// F

		// E

		win.IMShow(hsvFrame)
		keyPressed := win.WaitKey(1)
		if keyPressed == 'q' {
			return
		}
	}
}

func (w *WandaVision) Start() {
	w.topCamera = acquisition.CreateTopCameraAcquisition(w.conf)
	w.topCamera.Start()

	w.ballNarrow = ball.NewNarrowHaesveBall(w.conf)

	go worker(w)
	// worker(w)

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
