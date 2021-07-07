package main

import (
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
var greenWin = gocv.NewWindow("GREENMASK")

func main() {
	frame := gocv.NewMat()
	defer frame.Close()

	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	grayFrame := gocv.NewMat()
	defer grayFrame.Close()

	greenMask := gocv.NewMat()
	defer greenMask.Close()

	config, _ := configuration.LoadStartupConfig()

	topCamera := acquisition.CreateTopCameraAcquisition(&config)
	topCamera.Start()

	erodeMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()

	// dilateMat := gocv.Ones(3, 3, gocv.MatTypeCV8UC1)
	// defer dilateMat.Close()

	// rad := 100

	// cRed := color.RGBA{255, 0, 255, 0}
	cBlue := color.RGBA{0, 0, 255, 0}
	_ = cBlue

	// frameRead := gocv.IMRead("/Users/hariangr/Downloads/untitled.png", gocv.IMReadAnyColor)

	// midX := 0
	// midY := 0

	greenUpper := gocv.NewScalar(87, 255, 255, 1)
	greenLower := gocv.NewScalar(51, 0, 0, 0)

	greenErode := gocv.Ones(11, 11, gocv.MatTypeCV8UC1)
	defer greenErode.Close()
	greenErode2 := gocv.Ones(9, 9, gocv.MatTypeCV8UC1)
	defer greenErode2.Close()

	greenDilate := gocv.Ones(17, 17, gocv.MatTypeCV8UC1)
	defer greenDilate.Close()

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	var pointVecs gocv.PointsVector
	defer pointVecs.Close()
	// converHullMask := gocv.NewMat()

	for {
		// frameRead.CopyTo(&frame)
		topCamera.Read(&frame)

		topCamera.ReadHSV(&hsvFrame)

		topCamera.ReadGray(&grayFrame)
		gocv.Threshold(grayFrame, &grayFrame, 253, 255, gocv.ThresholdBinary)
		gocv.Erode(grayFrame, &grayFrame, erodeMat)

		gocv.InRangeWithScalar(hsvFrame, greenLower, greenUpper, &greenMask)
		gocv.Erode(greenMask, &greenMask, greenErode)
		gocv.Dilate(greenMask, &greenMask, greenDilate)
		gocv.Erode(greenMask, &greenMask, greenErode2)
		pointVecs = gocv.FindContoursWithParams(greenMask, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)
		// for i := 0; i < pointVecs.Size(); i++ {
		// 	it := pointVecs.At(i)
		// 	gocv.ConvexHull(it, &converHullMask, true, false)
		// }
		gocv.DrawContours(&frame, pointVecs, -1, cBlue, 3)

		win.IMShow(grayFrame)
		winRaw.IMShow(frame)
		greenWin.IMShow(greenMask)

		win.WaitKey(1)
	}

}
