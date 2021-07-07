package main

import (
	"image"
	"image/color"
	"math"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
	// "harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	// "harianugrah.com/brainfreeze/pkg/models/configuration"
)

var winHSV = gocv.NewWindow("HSV")
var winRaw = gocv.NewWindow("RAW")
var winGray = gocv.NewWindow("Gray")

func main() {
	frame := gocv.NewMat()
	defer frame.Close()

	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	grayFrame := gocv.NewMat()
	defer grayFrame.Close()

	config, _ := configuration.LoadStartupConfig()

	topCamera := acquisition.CreateTopCameraAcquisition(&config)
	topCamera.Start()

	// whiteLine := gocv.NewMat()

	// upper := gocv.NewScalar(0, 255, 255, 1)
	// lower := gocv.NewScalar(0, 0, 163, 0)

	erodeMat := gocv.Ones(9, 9, gocv.MatTypeCV8UC1)
	defer erodeMat.Close()

	// toShow := gocv.NewMat()
	// blackAll := gocv.Ones(640, 640, gocv.MatTypeCV8SC3)
	rad := float64(300)

	upperRed := gocv.NewScalar(179, 255, 255, 1)
	lowerRed := gocv.NewScalar(165, 56, 167, 0)

	hierarchyMat := gocv.NewMat()
	defer hierarchyMat.Close()

	var pointVecs gocv.PointsVector
	defer pointVecs.Close()

	// cRed := color.RGBA{255, 0, 255, 0}
	cBlue := color.RGBA{0, 255, 255, 0}

	for {
		topCamera.Read(&frame)

		topCamera.ReadGray(&grayFrame)
		gocv.Threshold(grayFrame, &grayFrame, 254, 255, gocv.ThresholdBinary)

		topCamera.ReadHSV(&hsvFrame)
		gocv.InRangeWithScalar(hsvFrame, lowerRed, upperRed, &hsvFrame)
		gocv.Erode(hsvFrame, &hsvFrame, erodeMat)
		pointVecs = gocv.FindContoursWithParams(hsvFrame, &hierarchyMat, gocv.RetrievalExternal, gocv.ChainApproxNone)

		detecteds := []models.DetectionObject{}
		for i := 0; i < pointVecs.Size(); i++ {
			it := pointVecs.At(i)
			area := gocv.ContourArea(it)

			if area < config.Wanda.MinimumHsvArea {
				// Skip kalau ukurannya kekecilan
				continue
			}

			rect := gocv.BoundingRect(it)
			// gocv.Rectangle(&frame, rect, cRed, 2)

			d := models.NewDetectionObject(rect)
			detecteds = append(detecteds, d)
		}

		lastOneWasWhite := true
		sudutPutih := make([]float64, 0)
		for i := -180.0; i <= 180.0; i++ {
			x := int(rad*math.Cos(i*math.Pi/180.0)) + 320
			y := int(rad*math.Sin(i*math.Pi/180.0)) + 320

			px := grayFrame.GetUCharAt(y, x)
			if px > 253 {
				if !lastOneWasWhite {
					gocv.Circle(&frame, image.Point{X: x, Y: y}, 5, color.RGBA{uint8(i + 90), 0, 0, 0}, -1)
					sudutPutih = append(sudutPutih, i)
					lastOneWasWhite = true
				}
			} else {
				lastOneWasWhite = false
			}
		}

		for i := 0; i < len(sudutPutih)-1; i++ {
			smaller := sudutPutih[i]
			bigger := sudutPutih[i+1]

			offsetsmaller := smaller + 180
			offsetsbigger := bigger + 180

			if offsetsbigger-offsetsmaller < 15 {
				continue
			}
			if offsetsbigger-offsetsmaller > 100 {
				continue
			}

			for _, v := range detecteds {
				sudut_merah := v.AsTransform(&config)
				if sudut_merah.RobROT < models.Degree(bigger) && sudut_merah.RobROT > models.Degree(smaller) {
					gocv.Line(&frame, image.Point{320, 320}, v.Midpoint, cBlue, 2)
				}
			}

		}

		winGray.IMShow(grayFrame)
		winHSV.IMShow(hsvFrame)
		winRaw.IMShow(frame)
		winHSV.WaitKey(1)
	}

}
