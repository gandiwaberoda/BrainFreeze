package main

import (
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

func main() {
	conf, _ := configuration.LoadStartupConfig()
	vc := acquisition.CreateTopCameraAcquisition(&conf)
	_ = vc
	vc.Start()

	win := gocv.NewWindow("apalah")

	for {
		frame := gocv.NewMat()
		vc.Read(&frame)

		win.IMShow(frame)
		win.WaitKey(1)
	}
}
