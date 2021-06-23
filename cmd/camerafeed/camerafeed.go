package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

func main() {
	vc, _ := gocv.VideoCaptureDevice(0)
	defer vc.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	vc.Read(&frame)

	fmt.Println("Rows:", frame.Rows())
	fmt.Println("Cols:", frame.Cols())
}
