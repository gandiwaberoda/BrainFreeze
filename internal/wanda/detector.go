package wanda

import (
	"fmt"
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
	res = models.SortDetectionsObjectByDistanceToPoint(image.Point{w.conf.Camera.PostWidth / 2, w.conf.Camera.PostHeight / 2}, res)

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

func detectDummyCircular(w *WandaVision, wg *sync.WaitGroup, tophsvFrame *gocv.Mat) {
	defer wg.Done()
	detecteds := w.dummyCircular.Detect(tophsvFrame)
	w.state.UpdateCircularDummy(detecteds)
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

func detectStraight(w *WandaVision, wg *sync.WaitGroup, forHsvFrame *gocv.Mat, forPostFrame *gocv.Mat) {
	defer wg.Done()
	detecteds := w.forStraight.Detect(forHsvFrame, forPostFrame)

	usedColorIds := make(map[int]bool)

	onceByClosest2Robot := make([]models.StraightDetection, 0)

	if len(detecteds) > 0 {
		for _, v := range detecteds {
			if _, exist := usedColorIds[v.DetectedColor.Id]; !exist {
				gocv.PutText(forPostFrame, "Entah", image.Pt(w.conf.Camera.ForMidX+40, v.LowerY), gocv.FontHersheyPlain, 1.3, v.DetectedColor.Visualize, 2)
				gocv.Line(forPostFrame, image.Pt(w.conf.Camera.ForMidX-20, v.LowerY), image.Pt(w.conf.Camera.ForMidX-20, v.UpperY), v.DetectedColor.Visualize, 5)
				onceByClosest2Robot = append(onceByClosest2Robot, v)

				usedColorIds[v.DetectedColor.Id] = true
			}
		}

		// obj := detecteds[0]

		// gocv.Line(hsvFrame, image.Point{320, 320}, obj.Midpoint, cBlue, 2)

		// t := obj.AsTransform(w.conf)
		// t.InjectWorldTransfromFromRobotTransform(w.state.GetState().MyTransform)
		// w.state.UpdateFriendGoalpostTransform(t)
	}

	simplified := make([]models.StraightDetectionObj, 0)
	for _, v := range onceByClosest2Robot {
		simplified = append(simplified, models.StraightDetectionObj{
			ClosestDistPx:     w.conf.Camera.PostHeight - v.LowerY,
			FurthestDistPx:    w.conf.Camera.PostHeight - v.UpperY,
			DetectedColorName: v.DetectedColor.Name,
		})
	}

	// fmt.Println(simplified)
	w.state.UpdateStraight(simplified)

	gocv.Circle(forPostFrame, image.Point{w.conf.Camera.ForMidX + 15, w.conf.Camera.ForPostHeight / 2}, 7, color.RGBA{255, 255, 255, 1}, 1)
	dist := fmt.Sprint(w.state.GetState().Araya.Dist0, "cm")
	gocv.PutText(forPostFrame, dist, image.Point{w.conf.Camera.ForMidX, w.conf.Camera.ForPostHeight / 2}, gocv.FontHersheyComplexSmall, 1.2, color.RGBA{255, 255, 255, 1}, 2)

}

// func detectMagenta(w *WandaVision, wg *sync.WaitGroup, topFrame *gocv.Mat, topHsvFrame *gocv.Mat) {

// 	wg.Done()
// }

func detectRadialAvoid(w *WandaVision, wg *sync.WaitGroup, topHsv *gocv.Mat) {
	defer wg.Done()
	detecteds := w.radialAvoid.Detect(topHsv)
	// w.state.UpdateCircularFieldLine(detecteds)
	_ = detecteds
}
