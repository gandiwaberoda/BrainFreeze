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
	"harianugrah.com/brainfreeze/internal/wanda/circular"
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

	fieldLineCircular *circular.FieldLineCircular
	goalpostCircular  *circular.GoalpostCircular

	magentaNarrow *magenta.NarrowHaesveMagenta
	forMagenta    *magenta.ForwardNarrowHaesveMagenta

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
	topHsvFrame := gocv.NewMat()
	defer topHsvFrame.Close()
	topGrayFrame := gocv.NewMat()
	defer topGrayFrame.Close()

	forFrame := gocv.NewMat()
	forHsvFrame := gocv.NewMat()
	defer forHsvFrame.Close()

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
		wg.Add(1)
		go detectForMagenta(w, &wg, &forHsvFrame)

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
		if found, result := w.ballNarrow.Detect(&forHsvFrame); found {
			if len(result) > 0 {
				for _, v := range result {
					gocv.Rectangle(&forFrame, v.Bbox, color.RGBA{255, 0, 0, 1}, 3)
				}
			}
		}

		// FPS Gauge
		fpsText := fmt.Sprint(w.fpsHsv.Read(), "FPS")
		gocv.PutText(&topFrame, fpsText, image.Point{10, 60}, gocv.FontHersheyPlain, 5, color.RGBA{0, 255, 255, 0}, 3)
		w.state.UpdateFpsHsv(w.fpsHsv.Read())

		wg.Wait()
		w.fpsHsv.Tick()

		if w.conf.Diagnostic.ShowScreen {
			// Put mid line to forward camera
			gocv.Line(&forHsvFrame, image.Point{w.conf.Camera.ForMidX, 0}, image.Point{w.conf.Camera.ForMidX, w.conf.Camera.ForPostHeight}, color.RGBA{255, 255, 255, 1}, 1)

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

			})
		}
	}
}

func (w *WandaVision) Start() {
	if hsvForwardWin == nil && w.conf.Diagnostic.ShowScreen {
		hsvWin = gocv.NewWindow("HSV")
		rawWin = gocv.NewWindow("Post Processed")
		grayWin = gocv.NewWindow("Gray Win Processed")

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

// Helper

func detectTopBall(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	defer wg.Done()
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

			// Kadang kadang ada bbox yang gak valid, skip aja
			if newer.Bbox.Min.X == 0 && newer.Bbox.Max.X == 640 && newer.BboxArea == 0 {
				return
			}
			obj := w.latestKnownBallDetection.Lerp(newer, w.conf.Wanda.LerpValue)

			transform := obj.AsTransform(w.conf)
			transform.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)

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
}

func detectMagenta(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	defer wg.Done()

	if w.conf.Wanda.DisableMagentaDetection {
		return
	}

	narrowMagentaFound, narrowMagentaRes := w.magentaNarrow.Detect(topHsvFrame)
	if narrowMagentaFound && len(narrowMagentaRes) > 0 {
		if !w.latestKnownMagentaSet {
			w.latestKnownMagentaDetection = narrowMagentaRes[0]
			w.latestKnownMagentaSet = true
		}
		sortedByDist := models.SortDetectionsObjectByDistanceToPoint(w.latestKnownMagentaDetection.Midpoint, narrowMagentaRes)
		newer := sortedByDist[0]
		obj := w.latestKnownMagentaDetection.Lerp(newer, w.conf.Wanda.LerpValue)
		w.latestKnownMagentaDetection = obj

		t := obj.AsTransform(w.conf)
		t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		w.state.UpdateMagentaTransform(t)

		// t := narrowMagentaRes[0].AsTransform(w.conf)
		// t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		// w.state.UpdateMagentaTransform(t)
		// gocv.Circle(topFrame, obj.Midpoint, obj.OuterRad, w.warnaNewest, 2)
	}
}

func detectCyan(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	defer wg.Done()

	if w.conf.Wanda.DisableCyanDetection {
		return
	}

	narrowCyanFound, narrowCyanRes := w.cyanNarrow.Detect(topHsvFrame)
	if narrowCyanFound && len(narrowCyanRes) > 0 {
		if !w.latestKnownCyanSet {
			w.latestKnownCyanDetection = narrowCyanRes[0]
			w.latestKnownCyanSet = true
		}
		sortedByDist := models.SortDetectionsObjectByDistanceToPoint(w.latestKnownCyanDetection.Midpoint, narrowCyanRes)
		newer := sortedByDist[0]
		obj := w.latestKnownCyanDetection.Lerp(newer, w.conf.Wanda.LerpValue)
		w.latestKnownCyanDetection = obj

		t := obj.AsTransform(w.conf)
		t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		w.state.UpdateCyanTransform(t)

		gocv.Line(topHsvFrame, image.Point{320, 320}, obj.Midpoint, color.RGBA{255, 0, 0, 1}, 2)
	}
}

func detectDummy(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {
	defer wg.Done()
	found, res := w.dummyNarrow.Detect(topHsvFrame)
	if !found {
		return
	}

	ts := []models.Transform{}
	for _, v := range res {
		t := v.AsTransform(w.conf)

		t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		ts = append(ts, t)
	}
	w.state.UpdateObstaclesTransform(ts)
}

func detectForGoalpost(w *WandaVision, wg *sync.WaitGroup, forFrame *gocv.Mat, forHsvFrame *gocv.Mat) {
	defer wg.Done()
	if found, result := w.goalpostHaesve.Detect(forHsvFrame); found {
		if len(result) > 0 {
			// result[0].
			sort.Slice(result, func(i, j int) bool {
				return result[i].BboxArea > result[j].BboxArea
			})

			gocv.Rectangle(forFrame, result[0].Bbox, w.warnaNewest, 3)
			gocv.PutText(forFrame, "Gawang", result[0].Bbox.Min, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 0, 255, 1}, 2)
		}

	}
}

func detectForMagenta(w *WandaVision, wg *sync.WaitGroup, forHsvFrame *gocv.Mat) {
	defer wg.Done()
	if found, result := w.forMagenta.Detect(forHsvFrame); found {
		if len(result) > 0 {
			sort.Slice(result, func(i, j int) bool {
				return result[i].ContourArea > result[j].ContourArea
			})

			v := result[0]
			gocv.Rectangle(forHsvFrame, v.Bbox, w.warnaNewest, 3)

			// Isi error line ke mid
			errorNotDegree := v.Midpoint.X - w.conf.Camera.ForMidX
			gocv.PutText(forHsvFrame, fmt.Sprint("Magenta: ", errorNotDegree), v.Midpoint, gocv.FontHersheyPlain, 1, color.RGBA{255, 255, 255, 0}, 1)
			gocv.Line(forHsvFrame, v.Midpoint, image.Point{w.conf.Camera.ForMidX, v.Midpoint.Y}, color.RGBA{128, 0, 0, 1}, 2)
		}

	}
}

func detectLineFieldCircular(w *WandaVision, wg *sync.WaitGroup, grayFrame *gocv.Mat) {
	defer wg.Done()
	detecteds := w.fieldLineCircular.Detect(grayFrame)
	w.state.UpdateCircularFieldLine(detecteds)
}

func detectGoalpostCircular(w *WandaVision, wg *sync.WaitGroup, hsvFrame *gocv.Mat, grayFrame *gocv.Mat) {
	cBlue := color.RGBA{0, 128, 255, 0}

	defer wg.Done()
	detecteds := w.goalpostCircular.Detect(hsvFrame, grayFrame)
	if len(detecteds) > 0 {
		obj := detecteds[0]

		gocv.Line(hsvFrame, image.Point{320, 320}, obj.Midpoint, cBlue, 2)

		t := obj.AsTransform(w.conf)
		t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		w.state.UpdateFriendGoalpostTransform(t)
	}
}

// func detectMagenta(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {

// 	wg.Done()
// }
