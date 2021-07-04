package wanda

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"sync"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/cyan"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/dummy"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/goalpost"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/magenta"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	"harianugrah.com/brainfreeze/pkg/models/state"
)

type WandaVision struct {
	isRunning      bool
	conf           *configuration.FreezeConfig
	topCamera      *acquisition.TopCameraAcquisition
	forwardCamera  *acquisition.ForwardCameraAcquisition
	ballNarrow     *ball.NarrowHaesveBall
	magentaNarrow  *magenta.NarrowHaesveMagenta
	dummyNarrow    *dummy.NarrowHaesveDummy
	goalpostHaesve *goalpost.HaesveGoalpost
	cyanNarrow     *cyan.NarrowHaesveCyan
	state          *state.StateAccess
	fpsHsv         *diagnostic.FpsGauge

	latestKnownBallDetection models.DetectionObject
	latestKnownBallSet       bool

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
// var hsvWin *gocv.Window
// var rawWin *gocv.Window

// var hsvForwardWin *gocv.Window
// var rawForwardWin *gocv.Window

var hsvWin = gocv.NewWindow("HSV")
var rawWin = gocv.NewWindow("Post Processed")

var hsvForwardWin = gocv.NewWindow("Forward HSV")
var rawForwardWin = gocv.NewWindow("Forward Post Processed")

func worker(w *WandaVision) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recovered from xx ", r)
			return
		}
	}()

	topFrame := gocv.NewMat()
	topHsvFrame := gocv.NewMat()
	defer topHsvFrame.Close()

	forFrame := gocv.NewMat()
	forHsvFrame := gocv.NewMat()
	defer forHsvFrame.Close()

	wg := sync.WaitGroup{}
	for {
		w.topCamera.Read(&topFrame)
		gocv.CvtColor(topFrame, &topHsvFrame, gocv.ColorBGRToHSV)
		gocv.GaussianBlur(topHsvFrame, &topHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

		w.forwardCamera.Read(&forFrame)
		gocv.CvtColor(forFrame, &forHsvFrame, gocv.ColorBGRToHSV)
		gocv.GaussianBlur(forHsvFrame, &forHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

		w.fpsHsv.Tick()

		// Ball
		wg.Add(1)
		go detectTopBall(w, &wg, &topFrame, &topHsvFrame)

		// EGP

		// FGP

		// Magenta
		wg.Add(1)
		go detectMagenta(w, &wg, &topFrame, &topHsvFrame)

		// Cyan
		wg.Add(1)
		go detectCyan(w, &wg, &topFrame, &topHsvFrame)

		// E

		// Dummy
		wg.Add(1)
		go detectDummy(w, &wg, &topFrame, &topHsvFrame)

		// Forward Goalpost
		wg.Add(1)
		go detectForGoalpost(w, &wg, &forFrame, &forHsvFrame)

		// Forward Ball
		if found, result := w.ballNarrow.Detect(&forHsvFrame); found {
			if len(result) > 0 {
				gocv.Rectangle(&topFrame, result[0].Bbox, w.warnaNewest, 3)
			}
		}

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&topFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		wg.Wait()

		if w.conf.Diagnostic.ShowScreen {
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
}

func (w *WandaVision) Start() {
	// if hsvForwardWin == nil && w.conf.Diagnostic.ShowScreen {
	// 	hsvWin = gocv.NewWindow("HSV")
	// 	rawWin = gocv.NewWindow("Post Processed")

	// 	hsvForwardWin = gocv.NewWindow("Forward HSV")
	// 	rawForwardWin = gocv.NewWindow("Forward Post Processed")
	// }

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

	w.topCenter = image.Point{
		w.conf.Camera.MidpointRad,
		w.conf.Camera.MidpointRad,
	}

	w.warnaNewest = color.RGBA{0, 255, 0, 0}
	w.warnaLastKnown = color.RGBA{0, 0, 255, 0}

	mainthread.Run(func() {
		worker(w)
	})

	fmt.Println("aw")

	w.isRunning = true
}

func (w *WandaVision) Stop() {
	w.topCamera.Stop()
}

// Helper

func detectTopBall(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	narrowBallFound, narrowBallRes := w.ballNarrow.Detect(topHsvFrame)
	if narrowBallFound {
		// TODO: Perlu lakukan classifier

		if len(narrowBallRes) > 0 {
			// newer := narrowBallRes[0]
			if !w.latestKnownBallSet {
				w.latestKnownBallDetection = narrowBallRes[0]
				w.latestKnownBallSet = true
			}
			sortedByDist := models.SortDetectionsObjectByDistanceToPoint(w.topCenter, narrowBallRes)
			newer := sortedByDist[0]
			obj := w.latestKnownBallDetection.Lerp(newer, w.conf.Wanda.LerpValue)

			transform := obj.AsTransform(w.conf)

			gocv.Rectangle(topFrame, narrowBallRes[0].Bbox, w.warnaNewest, 3)
			gocv.Circle(topFrame, obj.Midpoint, obj.OuterRad, w.warnaNewest, 2)
			// Origin to Ball Line
			gocv.Line(topFrame, w.conf.Camera.Midpoint, obj.Midpoint, w.warnaNewest, 2)

			w.state.UpdateBallTransform(transform)
			w.latestKnownBallDetection = obj
		} else {
			gocv.Line(topFrame, w.conf.Camera.Midpoint, w.latestKnownBallDetection.Midpoint, w.warnaLastKnown, 2)
		}
	} else {
		// Pake yang wide ball
		fmt.Println("loss")
	}
	wg.Done()
}

func detectMagenta(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	narrowMagentaFound, narrowMagentaRes := w.magentaNarrow.Detect(topHsvFrame)
	if narrowMagentaFound && len(narrowMagentaRes) > 0 {
		w.state.UpdateMagentaTransform(narrowMagentaRes[0].AsTransform(w.conf))
	}
	wg.Done()
}

func detectCyan(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	narrowCyanFound, narrowCyanRes := w.cyanNarrow.Detect(topHsvFrame)
	if narrowCyanFound && len(narrowCyanRes) > 0 {
		w.state.UpdateCyanTransform(narrowCyanRes[0].AsTransform(w.conf))
	}
	wg.Done()
}

func detectDummy(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	w.dummyNarrow.Detect(topHsvFrame)
	wg.Done()
}

func detectForGoalpost(w *WandaVision, wg *sync.WaitGroup, forFrame *gocv.Mat, forHsvFrame *gocv.Mat) {
	if found, result := w.goalpostHaesve.Detect(forHsvFrame); found {
		if len(result) > 0 {
			// result[0].
			sort.Slice(result, func(i, j int) bool {
				return result[i].BboxArea > result[j].BboxArea
			})

			gocv.Rectangle(forFrame, result[0].Bbox, w.warnaNewest, 3)
			gocv.PutText(forFrame, "Gawang", result[0].Bbox.Min, gocv.FontHersheyPlain, 1.2, w.warnaNewest, 2)
		}

	}
	wg.Done()
}

// func detectMagenta(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {

// 	wg.Done()
// }
