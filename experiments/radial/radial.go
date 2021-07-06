package main

import (
	"fmt"
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	// "harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	// "harianugrah.com/brainfreeze/pkg/models/configuration"
)

func calcColor(color int) (red, green, blue, alpha int) {
	alpha = color & 0xFF
	blue = (color >> 8) & 0xFF
	green = (color >> 16) & 0xFF
	red = (color >> 24) & 0xFF

	return red, green, blue, alpha
}

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

	erodeMat := gocv.Ones(1, 1, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()

	dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	defer dilateMat.Close()

	// rad := 100

	// cRed := color.RGBA{255, 0, 255, 0}
	cBlue := color.RGBA{0, 0, 255, 0}
	_ = cBlue

	frameRead := gocv.IMRead("/Users/hariangr/Downloads/untitled.png", gocv.IMReadAnyColor)

	// midX := 0
	// midY := 0

	for {
		frameRead.CopyTo(&frame)

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToGray)
		gocv.Threshold(hsvFrame, &hsvFrame, 250, 255, gocv.ThresholdBinary)

		fmt.Println(hsvFrame.Type())
		// gocv.Line(&frame, image.Point{midX, midY}, image.Point{midX + rad, midY}, cRed, 2)

		// for y := 0; y < frame.Rows(); y++ {
		// 	for x := 0; x < frame.Cols(); x++ {
		// 		px := hsvFrame.GetIntAt(y, x)

		// 		r, g, b, a := calcColor(px)

		// 		gocv.Circle(&frame)
		// 	}
		// }

		win.IMShow(hsvFrame)
		winRaw.IMShow(frame)
		win.WaitKey(1)
	}

}
