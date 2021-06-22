package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"
	"time"

	"gocv.io/x/gocv"
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

	src := config.Camera.Src[0]
	var errVc error
	var topCam *gocv.VideoCapture

	if len(src) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(src)
		if errInt != nil {
			panic(errInt)
		}
		topCam, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		topCam, errVc = gocv.VideoCaptureFile(config.Camera.Src[0])
	}

	if errVc != nil {
		panic(errVc)
	}

	firstFrame := gocv.NewMat()
	topCam.Read(&firstFrame)

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
		frame := gocv.NewMat()
		topCam.Read(&frame)

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
