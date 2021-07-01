package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
)

func main() {
	vc1, err1 := gocv.VideoCaptureDevice(0)

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

	erodeMat := gocv.Ones(27, 27, gocv.MatTypeCV8UC1)
	dilateMat := gocv.Ones(40, 40, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()
	defer dilateMat.Close()

	hsv := gocv.NewMat()
	for {
		gocv.CvtColor(f1, &hsv, gocv.ColorBGRToHSV)
		gocv.InRangeWithScalar(hsv, lower, upper, &hsv)

		gocv.Erode(hsv, &hsv, erodeMat)
		gocv.Dilate(hsv, &hsv, dilateMat)
		// // gocv.Threshold(gray, &gray, 170, 255, gocv.ThresholdBinary)

		if !f1.Empty() {
			win1.IMShow(f1)
			win2.IMShow(hsv)

			win1.WaitKey(1)
			win2.WaitKey(1)
		}
	}

}
