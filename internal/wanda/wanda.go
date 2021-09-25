package wanda

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/internal/wanda/circular"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/cyan"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/dummy"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/goalpost"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/magenta"
	"harianugrah.com/brainfreeze/internal/wanda/radial"
	"harianugrah.com/brainfreeze/internal/wanda/straight"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type WandaVision struct {
	isRunning bool
	conf      *configuration.FreezeConfig

	topCamera     *acquisition.TopCameraAcquisition
	forwardCamera *acquisition.ForwardCameraAcquisition
	fpsHsv        *diagnostic.FpsGauge
	state         *state.StateAccess

	ballNarrow     *ball.NarrowHaesveBall
	dummyNarrow    *dummy.NarrowHaesveDummy
	goalpostHaesve *goalpost.HaesveGoalpost
	cyanNarrow     *cyan.NarrowHaesveCyan

	radialAvoid *radial.TopRadialDummy

	fieldLineCircular *circular.FieldLineCircular
	goalpostCircular  *circular.GoalpostCircular

	magentaNarrow *magenta.NarrowHaesveMagenta
	forMagenta    *magenta.ForwardNarrowHaesveMagenta
	forStraight   *straight.ForwardColorStraight

	latestKnownBallDetection models.DetectionObject
	latestKnownBallSet       bool

	latestKnownCyanDetection    models.DetectionObject
	latestKnownCyanSet          bool
	latestKnownMagentaDetection models.DetectionObject
	latestKnownMagentaSet       bool

	topCenter image.Point

	warnaNewest    color.RGBA
	warnaLastKnown color.RGBA
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
var hsvWin *gocv.Window
var rawWin *gocv.Window
var grayWin *gocv.Window
var tempWin *gocv.Window

var hsvForwardWin *gocv.Window
var rawForwardWin *gocv.Window

func worker(w *WandaVision) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from xx ", r)
			return
		}
	}()

	topFrame := gocv.NewMat()
	defer topFrame.Close()
	topHsvFrame := gocv.NewMat()
	defer topHsvFrame.Close()
	topGrayFrame := gocv.NewMat()
	defer topGrayFrame.Close()

	forFrame := gocv.NewMat()
	defer forFrame.Close()
	forHsvFrame := gocv.NewMat()
	defer forHsvFrame.Close()
	vvv := gocv.NewMat()
	defer vvv.Close()

	wg := sync.WaitGroup{}
	for {
		w.topCamera.Read(&topFrame)
		w.topCamera.ReadHSV(&topHsvFrame)
		w.topCamera.ReadGray(&topGrayFrame)
		w.forwardCamera.Read(&forFrame)
		w.forwardCamera.ReadHSV(&forHsvFrame)

		// Ball
		wg.Add(1)
		go detectTopBall(w, &wg, &topFrame, &topHsvFrame)

		// EGP

		// FGP

		// Magenta
		wg.Add(1)
		go detectMagenta(w, &wg, &topFrame, &topHsvFrame)
		// wg.Add(1)
		// go detectForMagenta(w, &wg, &forHsvFrame)

		// Cyan
		wg.Add(1)
		go detectCyan(w, &wg, &topFrame, &topHsvFrame)

		// E

		// Dummy
		wg.Add(1)
		go detectDummy(w, &wg, &topFrame, &topHsvFrame)

		// // Forward Goalpost
		// wg.Add(1)
		// go detectForGoalpost(w, &wg, &forFrame, &forHsvFrame)

		// Circular Line Field
		wg.Add(1)
		go detectLineFieldCircular(w, &wg, &topGrayFrame)

		// Circular Goalpost
		wg.Add(1)
		go detectGoalpostCircular(w, &wg, &topHsvFrame, &topGrayFrame)

		// Forward Ball
		// if found, result := w.ballNarrow.Detect(&forHsvFrame); found {
		// 	if len(result) > 0 {
		// 		for _, v := range result {
		// 			gocv.Rectangle(&forFrame, v.Bbox, color.RGBA{255, 0, 0, 1}, 3)
		// 		}
		// 	}
		// }

		// zzz
		wg.Add(1)
		go detectStraight(w, &wg, &forHsvFrame, &forFrame)

		// wg.Add(1)
		// go detectRadialAvoid(w, &wg, &topHsvFrame)

		// gftt(&topGrayFrame)
		// hl(&topGrayFrame)
		whitte(&topFrame, &vvv)

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&topFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		wg.Wait()
		w.fpsHsv.Tick()

		if w.conf.Diagnostic.ShowScreen {
			// Put mid line to forward camera
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

				grayWin.IMShow(topGrayFrame)
				if keyPressed := grayWin.WaitKey(1); keyPressed == 'q' {
					return
				}

				tempWin.IMShow(vvv)
				if keyPressed := tempWin.WaitKey(1); keyPressed == 'q' {
					return
				}

			})
		}
	}
}

func (w *WandaVision) Start() {
	if hsvForwardWin == nil && w.conf.Diagnostic.ShowScreen {
		hsvWin = gocv.NewWindow("HSV")
		rawWin = gocv.NewWindow("Post Processed")
		grayWin = gocv.NewWindow("Gray Win Processed")
		tempWin = gocv.NewWindow("Entahlah")

		hsvForwardWin = gocv.NewWindow("Forward HSV")
		rawForwardWin = gocv.NewWindow("Forward Post Processed")
	}

	w.topCamera = acquisition.CreateTopCameraAcquisition(w.conf)
	w.topCamera.Start()

	w.forwardCamera = acquisition.NewForwardCameraAcquisition(w.conf)
	w.forwardCamera.Start()

	w.fpsHsv.Start()

	w.ballNarrow = ball.NewNarrowHaesveBall(w.conf)
	w.magentaNarrow = magenta.NewNarrowHaesveMagenta(w.conf)
	w.dummyNarrow = dummy.NewNarrowHaesveDummy(w.conf)
	w.cyanNarrow = cyan.NewNarrowHaesveCyan(w.conf)
	w.goalpostHaesve = goalpost.NewHaesveGoalpost(w.conf)
	w.fieldLineCircular = circular.NewFieldLineCircular(w.conf)
	w.goalpostCircular = circular.NewGoalpostCircular(w.conf)
	w.forMagenta = magenta.NewForwardNarrowHaesveMagenta(w.conf)
	w.forStraight = straight.NewForwardColorStraight(w.conf)
	w.radialAvoid = radial.NewRadialDummy(w.conf)

	w.topCenter = image.Point{
		w.conf.Camera.MidpointRad,
		w.conf.Camera.MidpointRad,
	}

	w.warnaNewest = color.RGBA{0, 255, 0, 0}
	w.warnaLastKnown = color.RGBA{0, 0, 255, 0}

	w.isRunning = true

	if w.conf.Diagnostic.ShowScreen {
		mainthread.Run(func() {
			worker(w)
		})
	} else {
		go worker(w)
	}
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}
