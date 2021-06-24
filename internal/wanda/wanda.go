package wanda

import (
	"fmt"
	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/dummy"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/magenta"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
	"image"
	"image/color"
)

type WandaVision struct {
	isRunning     bool
	conf          *configuration.FreezeConfig
	topCamera     *acquisition.TopCameraAcquisition
	ballNarrow    *ball.NarrowHaesveBall
	magentaNarrow *magenta.NarrowHaesveMagenta
	dummyNarrow   *dummy.NarrowHaesveDummy
	state         *state.StateAccess
	fpsHsv        *diagnostic.FpsGauge
}

func NewWandaVision(conf *configuration.FreezeConfig, state *state.StateAccess) *WandaVision {
	return &WandaVision{
		conf:   conf,
		state:  state,
		fpsHsv: diagnostic.NewFpsGauge(),
	}
}

// TODO: SIGSEGV Handle
// Penyebabnya karena Dilate

// Harus diluar mainthread
var hsvWin = gocv.NewWindow("HSV")
var rawWin = gocv.NewWindow("Post Processed")

func worker(w *WandaVision) {
	topCenter := image.Point{
		w.conf.Camera.MidpointRad,
		w.conf.Camera.MidpointRad,
	}

	warnaNewest := color.RGBA{0, 255, 0, 0}
	warnaLastKnown := color.RGBA{0, 0, 255, 0}

	frame := gocv.NewMat()
	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	var latestKnownBallDetection models.DetectionObject
	latestKnownBallSet := false

	for {
		w.topCamera.Read(&frame)
		w.fpsHsv.Tick()

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)
		// Blur
		gocv.GaussianBlur(hsvFrame, &hsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

		// Ball
		narrowBallFound, narrowBallRes := w.ballNarrow.Detect(&hsvFrame)
		if narrowBallFound {
			// TODO: Perlu lakukan classifier

			if len(narrowBallRes) > 0 {
				// newer := narrowBallRes[0]
				if !latestKnownBallSet {
					latestKnownBallDetection = narrowBallRes[0]
					latestKnownBallSet = true
				}
				sortedByDist := models.SortDetectionsObjectByDistanceToPoint(topCenter, narrowBallRes)
				newer := sortedByDist[0]
				obj := latestKnownBallDetection.Lerp(newer, w.conf.Wanda.LerpValue)

				transform := obj.AsTransform(w.conf)

				gocv.Rectangle(&frame, narrowBallRes[0].Bbox, warnaNewest, 3)
				gocv.Circle(&frame, obj.Midpoint, obj.OuterRad, warnaNewest, 2)
				// Origin to Ball Line
				gocv.Line(&frame, w.conf.Camera.Midpoint, obj.Midpoint, warnaNewest, 2)

				w.state.UpdateBallTransform(transform)
				latestKnownBallDetection = obj
			} else {
				gocv.Line(&frame, w.conf.Camera.Midpoint, latestKnownBallDetection.Midpoint, warnaLastKnown, 2)
			}
		} else {
			// Pake yang wide ball
			fmt.Println("loss")
		}

		// EGP

		// FGP

		// Magenta
		w.magentaNarrow.Detect(&hsvFrame)

		// E

		// Dummy
		w.dummyNarrow.Detect(&hsvFrame)

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&hsvFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		mainthread.Call(func() {
			rawWin.IMShow(frame)
			if keyPressed := rawWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			hsvWin.IMShow(hsvFrame)
			if keyPressed := hsvWin.WaitKey(1); keyPressed == 'q' {
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
	w.magentaNarrow = magenta.NewNarrowHaesveBall(w.conf)
	w.dummyNarrow = dummy.NewNarrowHaesveDummy(w.conf)

	mainthread.Run(func() {
		worker(w)
	})

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
