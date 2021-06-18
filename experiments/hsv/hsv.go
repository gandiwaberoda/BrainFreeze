package main

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/haesve/ball"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

func main() {
	vc, err := gocv.VideoCaptureFile("/Users/hariangr/Documents/MyFiles/Developer/Robotec/Beroda/ng/BrainDead/assets/captured/2021-04-23 14:50:34.157190 x 60.00024.1280.0.720.0.mp4")
	if err != nil {
		panic(err)
	}

	frame := gocv.NewMat()
	defer frame.Close()

	hsvFrame := gocv.NewMat()
	defer hsvFrame.Close()

	config, _ := configuration.LoadStartupConfig()
	ballHsv := ball.NewNarrowHaesveBall(&config)

	win := gocv.NewWindow("entah")

	circleMask := gocv.NewMatWithSize(720, 1280, gocv.MatTypeCV8UC3)
	mid := image.Point{
		X: 582, Y: 406,
	}
	white := color.RGBA{255, 255, 255, 0}
	gocv.Circle(&circleMask, mid, 287, white, -1)

	for {
		vc.Read(&frame)

		gocv.BitwiseAnd(frame, circleMask, &frame)

		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)
		ballHsv.Detect(hsvFrame)
		win.IMShow(frame)
		win.WaitKey(1)
	}

}
