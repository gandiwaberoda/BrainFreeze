package main

import (
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
	ballHsv := ball.CreateNewNarrowHaesveBall(config)

	for {
		vc.Read(&frame)
		gocv.CvtColor(frame, &hsvFrame, gocv.ColorBGRToHSV)
		ballHsv.Detect(hsvFrame)
	}

}
