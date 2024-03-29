package main

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/internal/diagnostic"
)

func main() {
	vc1, err1 := gocv.VideoCaptureDevice(1)
	vc2, err2 := gocv.VideoCaptureDevice(0)

	if err1 != nil {
		panic(err1)
	}

	if err2 != nil {
		panic(err2)
	}

	f1 := gocv.NewMat()
	defer f1.Close()

	f2 := gocv.NewMat()
	defer f2.Close()

	win1 := gocv.NewWindow("F1")
	win2 := gocv.NewWindow("F2")

	fps1 := diagnostic.NewFpsGauge()
	fps1.Start()
	fps2 := diagnostic.NewFpsGauge()
	fps2.Start()

	white := color.RGBA{128, 255, 0, 0}

	mutex1 := sync.Mutex{}
	mutex2 := sync.Mutex{}

	fmt.Println("A")
	go func() {
		for {
			mutex1.Lock()
			vc1.Read(&f1)
			gocv.PutText(&f1, "FPS A: "+strconv.Itoa(fps1.Read()), image.Point{10, 50}, gocv.FontHersheyPlain, 3, white, 2)
			fps1.Tick()
			mutex1.Unlock()
		}
	}()
	fmt.Println("B")
	go func() {
		for {
			mutex2.Lock()
			vc2.Read(&f2)
			gocv.PutText(&f2, "FPS B: "+strconv.Itoa(fps2.Read()), image.Point{10, 50}, gocv.FontHersheyPlain, 3, white, 2)
			fps2.Tick()
			mutex2.Unlock()
		}
	}()

	fmt.Println("C")
	// <-time.After(time.Second * 2)

	for {
		if !f1.Empty() && !f2.Empty() {
			fmt.Println("YES")
			break
		} else {
			if f1.Empty() {
				fmt.Println("F1 Empty")
			}
			if !f2.Empty() {
				fmt.Println("F2 Empty")
			}
		}
	}

	fmt.Println("D")
	for {
		// mutex1.Lock()
		// mutex2.Lock()

		if !f1.Empty() {
			win1.IMShow(f1)
			win1.WaitKey(1)
		}

		if !f2.Empty() {
			win2.IMShow(f2)
			win2.WaitKey(1)
		}

		// mutex1.Unlock()
		// mutex2.Unlock()
	}

}
