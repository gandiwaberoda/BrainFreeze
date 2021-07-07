package main

import (
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

	upper := gocv.NewScalar(179, 255, 255, 1)
	lower := gocv.NewScalar(165, 56, 167, 0)

	erodeMat := gocv.Ones(9, 9, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()

	dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	defer dilateMat.Close()

	rad := float64(315)

	cGreen := color.RGBA{0, 255, 255, 0}

	for {
		topCamera.Read(&frame)
		topCamera.ReadHSV(&hsvFrame)
		gocv.InRangeWithScalar(hsvFrame, lower, upper, &hsvFrame)
		gocv.Erode(hsvFrame, &hsvFrame, erodeMat)

		// lastOneWasWhite := true
		for i := -170.0; i <= 170.0; i += 3 {
			x := int(rad*math.Cos(i*math.Pi/180.0)) + 320
			y := int(rad*math.Sin(i*math.Pi/180.0)) + 320

			gocv.Line(&frame, image.Point{320, 320}, image.Point{x, y}, cGreen, 2)
			// x := 320
			// y := 320

			// fmt.Println(x, y)

			// px := hsvFrame.GetUCharAt(y, x)
			// // if math.Abs(px) > 0 {
			// if px > 253 {
			// 	// fmt.Println(px)
			// 	if !lastOneWasWhite {
			// 		gocv.Circle(&frame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
			// 		lastOneWasWhite = true
			// 	}
			// } else {
			// 	lastOneWasWhite = false
			// }
		}

		win.IMShow(hsvFrame)
		winRaw.IMShow(frame)
		win.WaitKey(1)
	}

}
