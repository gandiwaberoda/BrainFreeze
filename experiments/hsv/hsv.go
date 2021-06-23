package main

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	// "harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	// "harianugrah.com/brainfreeze/pkg/models/configuration"
)

func main() {
	// vc, err := gocv.VideoCaptureFile("/Users/hariangr/Documents/MyFiles/Developer/Robotec/Beroda/ng/BrainDead/assets/captured/2021-04-23 14:50:34.157190 x 60.00024.1280.0.720.0.mp4")

	vc, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		panic(err)
	}

	frame := gocv.NewMat()
	defer frame.Close()

	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	// config, _ := configuration.LoadStartupConfig()
	// ballHsv := ball.NewNarrowHaesveBall(&config)

	win := gocv.NewWindow("entah")

	mid := image.Point{
		X: 582, Y: 406,
	}
	white := color.RGBA{255, 255, 255, 0}

	for {
		vc.Read(&frame)
		circleMask := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)
		gocv.Circle(&circleMask, mid, 287, white, -1)
		maskedFrame := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)

		frame.CopyToWithMask(&maskedFrame, circleMask)

		// n := gocv.NewMat()

		// gocv.BitwiseAnd(frame, circleMask, &n)
		// frame.CopyToWithMask(&n, circleMask)

		// gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)
		// x := ballHsv.Detect(&hsvFrame)
		// fmt.Println(x)
		win.IMShow(maskedFrame)
		win.WaitKey(1)
	}

}
