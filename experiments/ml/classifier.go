package main

import (
	"fmt"
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

func main() {
	vc, _ := gocv.VideoCaptureDevice(0)
	defer vc.Close()

	backend := gocv.NetBackendDefault
	target := gocv.NetTargetCPU
	net := gocv.ReadNetFromTensorflow("assets/converted_savedmodel/model.savedmodel/saved_model.pb")
	if net.Empty() {
		panic("net empty")
	}
	defer net.Close()
	net.SetPreferableBackend(backend)
	net.SetPreferableTarget(target)

	win := gocv.NewWindow("Entah")
	defer win.Close()

	frame := gocv.NewMat()
	defer frame.Close()

	statusColor := color.RGBA{0, 255, 0, 0}

	for {
		vc.Read(&frame)

		blob := gocv.BlobFromImage(frame, 1.0, image.Pt(224, 224), gocv.NewScalar(0, 0, 0, 0), true, false)
		net.SetInput(blob, "input")
		prob := net.Forward("softmax2")
		probMat := prob.Reshape(1, 1)
		_, maxVal, _, maxLoc := gocv.MinMaxLoc(probMat)
		_ = maxLoc

		// desc := "Unknown"
		// if maxLoc.X < 1000 {
		// 	desc = descriptions[maxLoc.X]
		// }

		// status = fmt.Sprintf("description: %v, maxVal: %v\n", desc, maxVal)
		gocv.PutText(&frame, fmt.Sprintf("%f", maxVal), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		win.IMShow(frame)
		if keyPressed := win.WaitKey(1); keyPressed == 'q' {
			break
		}
	}
}
