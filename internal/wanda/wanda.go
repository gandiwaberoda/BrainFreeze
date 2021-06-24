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
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type WandaVision struct {
	isRunning     bool
	conf          *configuration.FreezeConfig
	topCamera     *acquisition.TopCameraAcquisition
	forwardCamera *acquisition.ForwardCameraAcquisition
	ballNarrow    *ball.NarrowHaesveBall
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

// Harus diluar mainthread
var hsvWin = gocv.NewWindow("HSV")
var rawWin = gocv.NewWindow("Post Processed")

var fHsvWin = gocv.NewWindow("Forward HSV")
var fRawWin = gocv.NewWindow("Forward Post Processed")

func worker(w *WandaVision) {
	topCenter := image.Point{
		w.conf.Camera.MidpointRad,
		w.conf.Camera.MidpointRad,
	}

	warnaNewest := color.RGBA{0, 255, 0, 0}
	warnaLastKnown := color.RGBA{0, 0, 255, 0}

	topFrame := gocv.NewMat()
	topHsv := gocv.NewMat()
	defer topHsv.Close()

	forwardFrame := gocv.NewMat()
	// forwardHsv := gocv.NewMat()
	defer forwardFrame.Close()

	var latestKnownBallDetection models.DetectionObject
	latestKnownBallSet := false

	for {
		w.fpsHsv.Tick()

		w.topCamera.Read(&topFrame)
		w.forwardCamera.Read(&forwardFrame)

		gocv.CvtColor(topFrame, &topHsv, gocv.ColorBGRToHSV)
		// gocv.CvtColor(forwardFrame, &forwardHsv, gocv.ColorBGRToHSV)

		// Blur
		gocv.GaussianBlur(topHsv, &topHsv, image.Point{7, 7}, 0, 0, gocv.BorderDefault)
		// gocv.GaussianBlur(forwardHsv, &forwardHsv, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

		// Ball
		narrowBallFound, narrowBallRes := w.ballNarrow.Detect(&topHsv)
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

				gocv.Rectangle(&topFrame, narrowBallRes[0].Bbox, warnaNewest, 3)
				gocv.Circle(&topFrame, obj.Midpoint, obj.OuterRad, warnaNewest, 2)
				// Origin to Ball Line
				gocv.Line(&topFrame, w.conf.Camera.Midpoint, obj.Midpoint, warnaNewest, 2)

				w.state.UpdateBallTransform(transform)
				latestKnownBallDetection = obj
			} else {
				gocv.Line(&topFrame, w.conf.Camera.Midpoint, latestKnownBallDetection.Midpoint, warnaLastKnown, 2)
			}

		} else {
			// Pake yang wide ball
			fmt.Println("loss")
		}

		// EGP

		// FGP

		// F

		// E

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&topHsv, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		mainthread.Call(func() {
			rawWin.IMShow(topFrame)
			if keyPressed := rawWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			hsvWin.IMShow(topHsv)
			if keyPressed := hsvWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			// fHsvWin.IMShow(forwardHsv)
			// if keyPressed := fHsvWin.WaitKey(1); keyPressed == 'q' {
			// 	return
			// }

			fRawWin.IMShow(forwardFrame)
			if keyPressed := fRawWin.WaitKey(1); keyPressed == 'q' {
				return
			}

		})
	}
}

func (w *WandaVision) Start() {
	w.topCamera = acquisition.CreateTopCameraAcquisition(w.conf)
	w.topCamera.Start()

	w.forwardCamera = acquisition.NewForwardCameraAcquisition(w.conf)
	w.forwardCamera.Start()

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
