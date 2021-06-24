package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

func main() {
	conf, _ := configuration.LoadStartupConfig()

	var errVc error
	var vc1, vc2 *gocv.VideoCapture

	srcA := conf.Camera.Src[0]
	if len(srcA) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(srcA)
		if errInt != nil {
			panic(errInt)
		}
		vc1, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		vc1, errVc = gocv.VideoCaptureFile(srcA)
	}
	if errVc != nil {
		panic(errVc)
	}

	srcB := conf.Camera.SrcForward[0]
	if len(srcA) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(srcB)
		if errInt != nil {
			panic(errInt)
		}
		vc2, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		vc2, errVc = gocv.VideoCaptureFile(srcB)
	}
	if errVc != nil {
		panic(errVc)
	}

	f1 := gocv.NewMat()
	// defer f1.Close()

	f2 := gocv.NewMat()
	// defer f2.Close()

	win1 := gocv.NewWindow("F1")
	win2 := gocv.NewWindow("F2")

	fps1 := diagnostic.NewFpsGauge()
	fps1.Start()
	fps2 := diagnostic.NewFpsGauge()
	fps2.Start()

	white := color.RGBA{255, 255, 255, 0}

	mutex1 := sync.Mutex{}
	mutex2 := sync.Mutex{}

	go func() {
		for {
			mutex1.Lock()
			vc1.Read(&f1)
			gocv.PutText(&f1, "FPS A: "+strconv.Itoa(fps1.Read()), image.Point{10, 50}, gocv.FontHersheyPlain, 3, white, 2)
			fps1.Tick()
			mutex1.Unlock()
		}
	}()
	go func() {
		for {
			mutex2.Lock()
			vc2.Read(&f2)
			gocv.PutText(&f2, "FPS B: "+strconv.Itoa(fps2.Read()), image.Point{10, 50}, gocv.FontHersheyPlain, 3, white, 2)
			fps2.Tick()
			mutex2.Unlock()
		}
	}()

	for {
		<-time.After(time.Second * 1)
		fmt.Println("Waiting camera ready...")
		if !f1.Empty() && !f2.Empty() {
			break
		}
	}

	for {
		// mutex1.Lock()
		// mutex2.Lock()

		win1.IMShow(f1)
		win2.IMShow(f2)

		win1.WaitKey(1)
		win2.WaitKey(1)

		// mutex1.Unlock()
		// mutex2.Unlock()
	}

}
