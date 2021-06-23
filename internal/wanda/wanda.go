package wanda

import (
	"fmt"
	"image"
	"image/color"

	// "github.com/faiface/mainthread"
	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
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
	fpsHsv     *diagnostic.FpsGauge
}

func NewWandaVision(conf *configuration.FreezeConfig, state *state.StateAccess) *WandaVision {
	return &WandaVision{
		conf:   conf,
		state:  state,
		fpsHsv: diagnostic.NewFpsGauge(),
	}
}

// Harus diluar mainthread
var win = gocv.NewWindow("hhh")

func worker(w *WandaVision) {
	warna := color.RGBA{0, 255, 0, 0}

	frame := gocv.NewMat()
	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	for {
		w.topCamera.Read(&frame)
		w.fpsHsv.Tick()

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)

		// Ball
		narrowBallRes := w.ballNarrow.Detect(&hsvFrame)
		if len(narrowBallRes) > 0 {
			transform := narrowBallRes[0].AsTransform(w.conf)
			gocv.Rectangle(&frame, narrowBallRes[0].Bbox, warna, 3)
			w.state.UpdateBallTransform(transform)
		} else if len(narrowBallRes) == 0 {
			// Pake yang wide ball
			fmt.Println("loss")
		}

		// EGP

		// FGP

		// F

		// E

		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&hsvFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)

		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		mainthread.Call(func() {

			win.IMShow(hsvFrame)
			keyPressed := win.WaitKey(1)
			if keyPressed == 'q' {
				return
			}

		})
	}
}

func (w *WandaVision) Start() {
	w.topCamera = acquisition.CreateTopCameraAcquisition(w.conf)
	w.topCamera.Start()

	w.fpsHsv.Start()

	w.ballNarrow = ball.NewNarrowHaesveBall(w.conf)

	mainthread.Run(func() {
		worker(w)
	})

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
