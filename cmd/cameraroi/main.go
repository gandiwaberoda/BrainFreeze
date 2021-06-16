package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/wanda/acquisition"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

var (
	tb_x   *gocv.Trackbar
	tb_y   *gocv.Trackbar
	tb_rad *gocv.Trackbar
)

func captioner() {
	ticker := time.NewTicker(time.Millisecond * 200)
	for {
		<-ticker.C
		fmt.Println("X:", tb_x.GetPos(), "\tY:", tb_y.GetPos(), "\tRAD:", tb_rad.GetPos())
	}
}

func main() {
	config, err := configuration.LoadStartupConfig()
	if err != nil {
		log.Fatalln("Gagal meload config", err)
	}

	topCam := acquisition.CreateTopCameraAcquisition(&config, false, false, false)
	topCam.Start()
	firstFrame := topCam.Read()

	win := gocv.NewWindow("Camera Region of Interest")

	wider := firstFrame.Cols()
	if firstFrame.Rows() > wider {
		wider = firstFrame.Rows()
	}
	tb_x = win.CreateTrackbar("X", firstFrame.Cols())
	tb_y = win.CreateTrackbar("Y", firstFrame.Rows())
	tb_rad = win.CreateTrackbar("Radius", wider)

	tb_x.SetPos(config.Camera.MidpointX)
	tb_y.SetPos(config.Camera.MidpointY)
	tb_rad.SetPos(config.Camera.MidpointRad)

	go captioner()

	for {
		frame := topCam.Read()
		mid := image.Point{
			X: tb_x.GetPos(),
			Y: tb_y.GetPos(),
		}
		color := color.RGBA{
			255, 0, 0, 0,
		}
		gocv.Circle(&frame, mid, tb_rad.GetPos(), color, 10)

		win.IMShow(frame)

		keyPressed := win.WaitKey(1)
		if keyPressed == 'q' {
			return
		}
	}
}
