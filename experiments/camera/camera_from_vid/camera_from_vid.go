package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"gocv.io/x/gocv"
)

func main() {
	log.Println("FPS camera test")

	webcam, err := gocv.VideoCaptureFile("/Users/hariangr/Documents/MyFiles/Developer/Robotec/Beroda/ng/assets/captured/2021-04-23 14:53:43.070299 x 60.00024.1280.0.720.0.mp4")
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
		elapsedUs := time.Since(startTime).Microseconds()
		fps := math.Pow(10.0, 6.0) / float64(elapsedUs)
		elapsedStr := "FPS: " + fmt.Sprintf("%f", fps)

		gocv.PutText(&img, elapsedStr, fpsPos, gocv.FontHersheyPlain, 1.5, fpsColor, 1)

		prevWindow.IMShow(img)

		keyPressed := prevWindow.WaitKey(1)
		if keyPressed == 'q' {
			break
		}
	}
}
