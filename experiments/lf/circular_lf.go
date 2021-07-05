package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	// "harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	// "harianugrah.com/brainfreeze/pkg/models/configuration"
)

var win = gocv.NewWindow("HSV")
var winRaw = gocv.NewWindow("RAW")

func main() {
	frame := gocv.NewMat()
	defer frame.Close()

	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	config, _ := configuration.LoadStartupConfig()

	topCamera := acquisition.CreateTopCameraAcquisition(&config)
	topCamera.Start()

	// whiteLine := gocv.NewMat()

	// upper := gocv.NewScalar(0, 255, 255, 1)
	// lower := gocv.NewScalar(0, 0, 163, 0)

	erodeMat := gocv.Ones(1, 1, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()

	dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	defer dilateMat.Close()

	// toShow := gocv.NewMat()
	// blackAll := gocv.Ones(640, 640, gocv.MatTypeCV8SC3)
	rad := float64(130)

	for {
		// blackAll.CopyTo(&toShow)

		topCamera.Read(&frame)
		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToGray)
		gocv.Threshold(hsvFrame, &hsvFrame, 254, 256, gocv.ThresholdBinary)
		// topCamera.ReadHSV(&hsvFrame)

		// gocv.InRangeWithScalar(hsvFrame, lower, upper, &hsvFrame)

		// gocv.Erode(hsvFrame, &hsvFrame, erodeMat)
		gocv.Dilate(hsvFrame, &hsvFrame, dilateMat)

		// str := ""

		// for y := 0; y < 640; y++ {
		// 	for x := 0; x < 640; x++ {
		// 		str += fmt.Sprint(hsvFrame.GetIntAt(x, y)) + " "
		// 	}
		// 	str += "\n"
		// 	fmt.Println(y)
		// }

		// _ = ioutil.WriteFile("./temp.txt", []byte(str), 0644)
		// panic("")

		// gocv.Canny(hsvFrame, &hsvFrame, 50, 200)
		// gocv.Dilate(hsvFrame, &hsvFrame, dilateMat)

		for i := -90.0; i <= 90.0; i++ {
			x := int(rad*math.Cos(i*math.Pi/180.0)) + 320
			y := int(rad*math.Sin(i*math.Pi/180.0)) + 320
			// x := 320
			// y := 320

			// fmt.Println(x, y)

			px := hsvFrame.GetIntAt(y, x)
			// if math.Abs(px) > 0 {
			if px != 0 {
				fmt.Println(px)
				gocv.Circle(&frame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
			}
		}

		win.IMShow(hsvFrame)
		winRaw.IMShow(frame)
		win.WaitKey(1)
	}

}
