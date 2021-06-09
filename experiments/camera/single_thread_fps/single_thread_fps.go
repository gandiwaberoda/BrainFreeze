package main

import (
	"image"
	"image/color"
	"log"
	"strconv"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	log.Println("FPS HSV camera test")

	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Fatalf("Failed to open video capture device")
	}
	defer webcam.Close()

	prevWindow := gocv.NewWindow("Preview Window")
	defer prevWindow.Close()

	img := gocv.NewMat()
	defer img.Close()

	for {
		startTime := time.Now()

		webcam.Read(&img)

		fpsPos := image.Point{X: 10, Y: 40}
		fpsColor := color.RGBA{255, 255, 255, 0}
		elapsed := time.Now().Sub(startTime)
		fps := 1000.0 / elapsed.Milliseconds()
		elapsedStr := "FPS: " + strconv.Itoa(int(fps))

		gocv.PutText(&img, elapsedStr, fpsPos, gocv.FontHersheyPlain, 1.5, fpsColor, 1)

		prevWindow.IMShow(img)

		keyPressed := prevWindow.WaitKey(1)
		if keyPressed == 'q' {
			break
		}
	}

}
