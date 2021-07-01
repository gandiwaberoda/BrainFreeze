package main

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"strconv"
	"sync"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
)

type CountourArea struct {
	Area    float64
	PVector gocv.PointVector
}

func main() {
	vc1, err1 := gocv.VideoCaptureDevice(1)

	if err1 != nil {
		panic(err1)
	}

	f1 := gocv.NewMat()
	defer f1.Close()

	win1 := gocv.NewWindow("F1")
	win2 := gocv.NewWindow("F2")

	fps1 := diagnostic.NewFpsGauge()
	fps1.Start()

	white := color.RGBA{128, 255, 0, 0}

	mutex1 := sync.Mutex{}

	fmt.Println("A")
	go func() {
		for {
			mutex1.Lock()
			vc1.Read(&f1)
			gocv.PutText(&f1, "FPS A: "+strconv.Itoa(fps1.Read()), image.Point{10, 50}, gocv.FontHersheyPlain, 3, white, 2)
			fps1.Tick()
			mutex1.Unlock()
		}
	}()

	fmt.Println("C")
	// <-time.After(time.Second * 2)

	for {
		if !f1.Empty() {
			fmt.Println("YES")
			break
		} else {
			if f1.Empty() {
				fmt.Println("F1 Empty")
			}
		}
	}

	fmt.Println("D")

	upper := gocv.NewScalar(166, 51, 255, 1)
	lower := gocv.NewScalar(0, 0, 244, 0)

	dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	erodeMat := gocv.Ones(35, 35, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()
	defer dilateMat.Close()

	hsv := gocv.NewMat()
	contur := gocv.NewMat()

	c := color.RGBA{0, 128, 128, 0}

	for {
		gocv.CvtColor(f1, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lower, upper, &hsv)

		gocv.Dilate(hsv, &hsv, dilateMat)
		gocv.Erode(hsv, &hsv, erodeMat)
		// // gocv.Threshold(gray, &gray, 170, 255, gocv.ThresholdBinary)

		// gocv.Canny(hsv, &hsv, 30, 10)
		// gocv.HoughLines(hsv, &hsv, math.Pi/180, 50, 20)
		// gocv.HoughLinesP(hsv, &hsv, math.Pi/180, 30, 10)

		pointVecs := gocv.FindContoursWithParams(hsv, &contur, gocv.RetrievalExternal, gocv.ChainApproxNone)
		rects := make([]CountourArea, 0)

		for i := 0; i < pointVecs.Size(); i++ {
			it := pointVecs.At(i)
			area := gocv.ContourArea(it)

			rects = append(rects, CountourArea{
				Area:    area,
				PVector: it,
			})

			// if area < n.conf.Wanda.MinimumHsvArea {
			// 	// Skip kalau ukurannya kekecilan
			// 	continue
			// }
			// if area > n.conf.Wanda.MaximumHsvArea {
			// 	continue
			// }

			// rect := gocv.BoundingRect(it)
			// gocv.Rectangle(&hsv, rect, c, 2)
			// gocv.PutText(&hsv, "Contur "+strconv.Itoa(i), rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)

			// d := models.NewDetectionObject(rect)
			// detecteds = append(detecteds, d)
		}

		sort.Slice(rects, func(i, j int) bool {
			return rects[i].Area > rects[j].Area
		})

		// for i := 0; i < pointVecs.Size(); i++ {
		// 	it := rects[i]

		// 	// if area < n.conf.Wanda.MinimumHsvArea {
		// 	// 	// Skip kalau ukurannya kekecilan
		// 	// 	continue
		// 	// }
		// 	// if area > n.conf.Wanda.MaximumHsvArea {
		// 	// 	continue
		// 	// }

		// 	rect := gocv.BoundingRect(it.PVector)
		// 	gocv.Rectangle(&hsv, rect, c, 2)
		// 	gocv.PutText(&hsv, "Contur "+strconv.Itoa(i), rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)

		// 	// d := models.NewDetectionObject(rect)
		// 	// detecteds = append(detecteds, d)
		// }

		if len(rects) > 0 {
			it := rects[0]
			rect := gocv.BoundingRect(it.PVector)
			gocv.Rectangle(&hsv, rect, c, 2)

			center := image.Point{
				X: (rect.Max.X + rect.Min.X) / 2,
				Y: (rect.Max.Y + rect.Min.Y) / 2,
			}

			gocv.Circle(&hsv, center, 10, c, -1)

			gocv.PutText(&hsv, "Gawang", rect.Min, gocv.FontHersheyPlain, 1.2, c, 2)
		}

		if !f1.Empty() {
			win1.IMShow(f1)
			win2.IMShow(hsv)

			win1.WaitKey(1)
			win2.WaitKey(1)
		}
	}

}
