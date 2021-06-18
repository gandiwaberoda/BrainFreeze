package wanda

import (
	"fmt"

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
	hsvFrame := gocv.NewMat()

	for {
		frame := w.topCamera.Read()

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)

		// Ball
		narrowBallRes := w.ballNarrow.Detect(hsvFrame)
		if len(narrowBallRes) > 0 {
			w.state.UpdateBallTransform(narrowBallRes[0].AsTransform())
		}

		// EGP

		// FGP

		// F

		// E
	}
}

func (w *WandaVision) Start() {
	w.topCamera = acquisition.CreateTopCameraAcquisition(w.conf)
	w.topCamera.Start()

	w.ballNarrow = ball.NewNarrowHaesveBall(w.conf)

	go worker(w)

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
