package acquisition

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock      *sync.Mutex
	IsRunning bool
	vc        *gocv.VideoCapture
	conf      *configuration.FreezeConfig
	Frame     gocv.Mat
	imgChan   chan gocv.Mat
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig) *TopCameraAcquisition {
	return &TopCameraAcquisition{
		IsRunning: false,
		conf:      conf,
		Lock:      &sync.Mutex{},
		imgChan:   make(chan gocv.Mat),
	}
}

func showImg(c *TopCameraAcquisition) {
	// now we can run stuff on the main thread like this
	mainthread.CallNonBlock(func() {
		fmt.Println("printing from the main thread")

		prevWindow := gocv.NewWindow("Preview Window")
		defer prevWindow.Close()
		mat := gocv.NewMat()

		for {
			mat = <-c.imgChan
			// prevWindow.IMShow(topCamera.Read())
			prevWindow.IMShow(mat)

			keyPressed := prevWindow.WaitKey(1)
			if keyPressed == 'q' {
				return
			}

		}

	})
	fmt.Println("printing from another thread")
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
		c.imgChan <- frame
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
	mainthread.Run(func() {
		showImg(c)
	})
}

func (c *TopCameraAcquisition) Stop() {
	c.vc.Close()
}
