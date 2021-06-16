package acquisition

import (
	// "image"
	// "image/color"
	// "math"
	// "time"
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock         *sync.RWMutex
	IsRunning    bool
	vc           *gocv.VideoCapture
	conf         *configuration.FreezeConfig
	Frame        gocv.Mat
	previewImage bool
	threaded     bool
	preprocess   bool
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig, threaded bool, preprocess bool, previewImage bool) *TopCameraAcquisition {

	return &TopCameraAcquisition{
		IsRunning:    false,
		conf:         conf,
		Lock:         &sync.RWMutex{},
		previewImage: previewImage,
		threaded:     threaded,
		// circleMask:   circleMask,
		Frame:      gocv.NewMat(),
		preprocess: preprocess,
	}
}

func worker(c *TopCameraAcquisition) {
	vvc, _ := gocv.VideoCaptureFile(c.conf.Camera.Src[0])
	defer vvc.Close()

	// frame := gocv.NewMatWithSize(c.conf.Camera.RawHeight, c.conf.Camera.RawWidth, gocv.MatTypeCV8UC3)
	frame := gocv.NewMat()

	circleMask := gocv.Zeros(c.conf.Camera.RawHeight, c.conf.Camera.RawWidth, gocv.MatTypeCV8UC1)
	mid := image.Point{
		X: c.conf.Camera.MidpointX, Y: c.conf.Camera.MidpointY,
	}
	white := color.RGBA{255, 255, 255, 0}
	gocv.Circle(&circleMask, mid, c.conf.Camera.MidpointRad, white, -1)

	for {
		// c.vc.Read(&frame)
		vvc.Read(&frame)

		// if c.preprocess {
		// 	preprocessTopCameraFrame(&frame)
		// }
		// m := gocv.Zeros(300, 300, gocv.MatTypeCV8UC1)
		// gocv.Circle()
		// masked := gocv.NewMat()
		// frame.CopyToWithMask(&masked, circleMask)

		c.Lock.Lock()
		c.Frame = circleMask
		c.Lock.Unlock()
	}
}

// func preprocessTopCameraFrame(frame *gocv.Mat) {

// 	startTime := time.Now()

// 	src := frame.Clone()

// 	fpsPos := image.Point{X: 10, Y: 40}
// 	fpsColor := color.RGBA{255, 255, 255, 0}
// 	elapsedUs := time.Since(startTime).Microseconds()
// 	fps := math.Pow(10.0, 6.0) / float64(elapsedUs)
// 	elapsedStr := "FPS: " + fmt.Sprintf("%f", fps)

// 	gocv.PutText(frame, elapsedStr, fpsPos, gocv.FontHersheyPlain, 1.5, fpsColor, 1)
// 	// frame = res.Clone()
// }

func (c *TopCameraAcquisition) Read() gocv.Mat {
	frame := gocv.NewMat()
	defer frame.Close()

	if !c.threaded {
		c.Lock.Lock()
		defer c.Lock.Unlock()

		c.vc.Read(&frame)
		c.Frame = frame.Clone()

		return c.Frame
	}

	c.Lock.RLock()
	defer c.Lock.RUnlock()
	return c.Frame
}

func (c *TopCameraAcquisition) Start() {
	src := c.conf.Camera.Src[0]

	var vc *gocv.VideoCapture
	var errVc error

	if len(src) == 1 {
		// Berupa angka
		srcInt, err := strconv.Atoi(src)
		if err != nil {
			panic(err)
		}
		vc, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Berupa video
		vc, errVc = gocv.VideoCaptureFile(src)
	}

	if errVc != nil {
		panic("can't open camera")
	}
	c.vc = vc

	if c.threaded {
		go worker(c)
	}

	if c.previewImage {
		mainthread.Run(func() {
			showImg(c)
		})
	}

}

func (c *TopCameraAcquisition) Stop() {
	c.vc.Close()
}

// func showImg(c *TopCameraAcquisition) {
// 	// now we can run stuff on the main thread like this
// 	mainthread.CallNonBlock(func() {
// 		fmt.Println("printing from the main thread")

// 		prevWindow := gocv.NewWindow("Preview Window")
// 		defer prevWindow.Close()
// 		// mat := gocv.NewMat()

// 		for {
// 			// mat = <-c.imgChan
// 			prevWindow.IMShow(c.Frame)

// 			keyPressed := prevWindow.WaitKey(1)
// 			if keyPressed == 'q' {
// 				return
// 			}

// 		}

// 	})
// 	fmt.Println("printing from another thread")
// }
