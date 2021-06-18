package wanda

import (
	"image"
	"image/color"
	"strconv"
	"time"

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

var frameCountSec int = 0
var fps = 0
var lastCheck = time.Now()

func worker(w *WandaVision) {
	win := gocv.NewWindow("hhh")

	warna := color.RGBA{0, 255, 0, 0}
	hsvFrame := gocv.NewMat()

	go func() {
		ticker := time.NewTicker(time.Millisecond * 1500)
		for {
			<-ticker.C

			elapsed := time.Since(lastCheck)

			fps = frameCountSec / int(elapsed.Seconds())

			frameCountSec = 0
			lastCheck = time.Now()
		}
	}()

	frame := gocv.NewMat()
	for {
		w.topCamera.Read(&frame)
		frameCountSec++

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

		// elapsed := time.Since(started)
		gocv.PutText(&hsvFrame, strconv.Itoa(fps), image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)

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
