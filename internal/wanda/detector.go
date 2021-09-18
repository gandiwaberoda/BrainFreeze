package wanda

import (
	"image"
	"image/color"
	"sync"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models"
)

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
		// fmt.Println("loss")
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

// func detectForGoalpost(w *WandaVision, wg *sync.WaitGroup, forFrame *gocv.Mat, forHsvFrame *gocv.Mat) {
// 	defer wg.Done()
// 	if found, result := w.goalpostHaesve.Detect(forHsvFrame); found {
// 		if len(result) > 0 {
// 			// result[0].
// 			sort.Slice(result, func(i, j int) bool {
// 				return result[i].BboxArea > result[j].BboxArea
// 			})

// 			gocv.Rectangle(forFrame, result[0].Bbox, w.warnaNewest, 3)
// 			gocv.PutText(forFrame, "Gawang", result[0].Bbox.Min, gocv.FontHersheyPlain, 1.2, color.RGBA{0, 0, 255, 1}, 2)
// 		}

// 	}
// }

// func detectForMagenta(w *WandaVision, wg *sync.WaitGroup, forHsvFrame *gocv.Mat) {
// 	defer wg.Done()
// 	if found, result := w.forMagenta.Detect(forHsvFrame); found {
// 		if len(result) > 0 {
// 			sort.Slice(result, func(i, j int) bool {
// 				return result[i].ContourArea > result[j].ContourArea
// 			})

// 			v := result[0]
// 			gocv.Rectangle(forHsvFrame, v.Bbox, w.warnaNewest, 3)

// 			// Isi error line ke mid
// 			errorNotDegree := v.Midpoint.X - w.conf.Camera.ForMidX
// 			gocv.PutText(forHsvFrame, fmt.Sprint("Magenta: ", errorNotDegree), v.Midpoint, gocv.FontHersheyPlain, 1, color.RGBA{255, 255, 255, 0}, 1)
// 			gocv.Line(forHsvFrame, v.Midpoint, image.Point{w.conf.Camera.ForMidX, v.Midpoint.Y}, color.RGBA{128, 0, 0, 1}, 2)
// 		}

// 	}
// }

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
