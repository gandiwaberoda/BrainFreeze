package main

import (
	"log"

	"gocv.io/x/gocv"
)

func main() {
	log.Println("Entah camera test")

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
		webcam.Read(&img)

		prevWindow.IMShow(img)

		keyPressed := prevWindow.WaitKey(1)
		if keyPressed == 'q' {
			break
		}
	}

}
