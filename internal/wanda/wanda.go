package wanda

import (
	"fmt"
	"image"
	"image/color"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/cyan"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/dummy"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/magenta"
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
	magentaNarrow *magenta.NarrowHaesveMagenta
	dummyNarrow   *dummy.NarrowHaesveDummy
	cyanNarrow    *cyan.NarrowHaesveCyan
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

var hsvForwardWin = gocv.NewWindow("Forward HSV")
var rawForwardWin = gocv.NewWindow("Forward Post Processed")

func worker(w *WandaVision) {
	topCenter := image.Point{
		w.conf.Camera.MidpointRad,
		w.conf.Camera.MidpointRad,
	}

	warnaNewest := color.RGBA{0, 255, 0, 0}
	warnaLastKnown := color.RGBA{0, 0, 255, 0}

	topFrame := gocv.NewMat()
	topHsvFrame := gocv.NewMat()
	defer topHsvFrame.Close()

	forFrame := gocv.NewMat()
	forHsvFrame := gocv.NewMat()
	defer forHsvFrame.Close()

	var latestKnownBallDetection models.DetectionObject
	latestKnownBallSet := false

	for {
		w.topCamera.Read(&topFrame)
		w.forwardCamera.Read(&forFrame)

		w.fpsHsv.Tick()

		// Ubah ke HSV
		gocv.CvtColor(topFrame, &topHsvFrame, gocv.ColorBGRToHSV)
		gocv.CvtColor(forFrame, &forHsvFrame, gocv.ColorBGRToHSV)

		// Blur
		gocv.GaussianBlur(topHsvFrame, &topHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)
		gocv.GaussianBlur(forHsvFrame, &forHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

		// Ball
		narrowBallFound, narrowBallRes := w.ballNarrow.Detect(&topHsvFrame)
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

		// Magenta
		narrowMagentaFound, narrowMagentaRes := w.magentaNarrow.Detect(&topHsvFrame)
		if narrowMagentaFound && len(narrowMagentaRes) > 0 {
			w.state.UpdateMagentaTransform(narrowMagentaRes[0].AsTransform(w.conf))
		}

		// Cyan
		narrowCyanFound, narrowCyanRes := w.cyanNarrow.Detect(&topHsvFrame)
		if narrowCyanFound && len(narrowCyanRes) > 0 {
			w.state.UpdateCyanTransform(narrowCyanRes[0].AsTransform(w.conf))
		}

		// E

		// Dummy
		w.dummyNarrow.Detect(&topHsvFrame)

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&topHsvFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		mainthread.Call(func() {
			rawWin.IMShow(topFrame)
			if keyPressed := rawWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			hsvWin.IMShow(topHsvFrame)
			if keyPressed := hsvWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			rawForwardWin.IMShow(forFrame)
			if keyPressed := rawForwardWin.WaitKey(1); keyPressed == 'q' {
				return
			}

			hsvForwardWin.IMShow(forHsvFrame)
			if keyPressed := hsvForwardWin.WaitKey(1); keyPressed == 'q' {
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
	w.magentaNarrow = magenta.NewNarrowHaesveMagenta(w.conf)
	w.dummyNarrow = dummy.NewNarrowHaesveDummy(w.conf)
	w.cyanNarrow = cyan.NewNarrowHaesveCyan(w.conf)

	mainthread.Run(func() {
		worker(w)
	})

	fmt.Println("aw")

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
