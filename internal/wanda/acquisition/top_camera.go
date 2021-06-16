package acquisition

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock      *sync.Mutex
	IsRunning bool
	vc        *gocv.VideoCapture
	conf      *configuration.FreezeConfig
	Frame     gocv.Mat
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig) *TopCameraAcquisition {
	return &TopCameraAcquisition{
		IsRunning: false,
		conf:      conf,
		Lock:      &sync.Mutex{},
	}
}

func worker(c *TopCameraAcquisition) {
	frame := gocv.NewMat()

	for {
		startTime := time.Now()

		c.vc.Read(&frame)

		fpsPos := image.Point{X: 10, Y: 40}
		fpsColor := color.RGBA{255, 255, 255, 0}
		elapsedUs := time.Since(startTime).Microseconds()
		fps := math.Pow(10.0, 6.0) / float64(elapsedUs)
		elapsedStr := "FPS: " + fmt.Sprintf("%f", fps)

		gocv.PutText(&frame, elapsedStr, fpsPos, gocv.FontHersheyPlain, 1.5, fpsColor, 1)

		c.Lock.Lock()
		c.Frame = frame
		c.Lock.Unlock()
	}
}

func PreprocessTopCameraFrame(frame *gocv.Mat) {

}

func (c *TopCameraAcquisition) Read() gocv.Mat {
	c.Lock.Lock()
	defer c.Lock.Unlock()
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

	go worker(c)
}

func (c *TopCameraAcquisition) Stop() {
	c.vc.Close()
}
