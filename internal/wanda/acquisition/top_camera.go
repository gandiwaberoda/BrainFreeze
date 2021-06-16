package acquisition

import (
	"image"
	"image/color"
	"strconv"
	"sync"

	"github.com/faiface/mainthread"
	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock  *sync.RWMutex
	vc    *gocv.VideoCapture
	conf  *configuration.FreezeConfig
	frame gocv.Mat
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig, threaded bool, preprocess bool, previewImage bool) *TopCameraAcquisition {
	return &TopCameraAcquisition{
		conf: conf,
		Lock: &sync.RWMutex{},
	}
}

func worker(c *TopCameraAcquisition) {
	for {
		c.Lock.Lock()
		c.read()
		c.Lock.Unlock()
	}
}

func (c *TopCameraAcquisition) read() {
	circleMask := gocv.Zeros(c.conf.Camera.RawHeight, c.conf.Camera.RawWidth, gocv.MatTypeCV8UC1)
	mid := image.Point{
		X: c.conf.Camera.MidpointX, Y: c.conf.Camera.MidpointY,
	}
	white := color.RGBA{255, 255, 255, 0}
	gocv.Circle(&circleMask, mid, c.conf.Camera.MidpointRad, white, -1)

	f := gocv.NewMatWithSize(c.conf.Camera.RawHeight, c.conf.Camera.RawWidth, gocv.MatTypeCV8UC3)
	c.vc.Read(&f)
	res := gocv.NewMat()
	f.CopyToWithMask(&res, circleMask)

	x0 := c.conf.Camera.MidpointX - c.conf.Camera.MidpointRad
	y0 := c.conf.Camera.MidpointY - c.conf.Camera.MidpointRad
	x1 := c.conf.Camera.MidpointX + c.conf.Camera.MidpointRad
	y1 := c.conf.Camera.MidpointY + c.conf.Camera.MidpointRad
	rect := image.Rect(x0, y0, x1, y1)

	c.frame = res.Region(rect)
}

func (c *TopCameraAcquisition) Read() gocv.Mat {
	return c.frame
}

func (c *TopCameraAcquisition) Start() {
	src := c.conf.Camera.Src[0]

	var vc *gocv.VideoCapture
	var errVc error
	if len(src) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(src)
		if errInt != nil {
			panic(errInt)
		}
		vc, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		vc, errVc = gocv.VideoCaptureFile(c.conf.Camera.Src[0])
	}
	if errVc != nil {
		panic(errVc)
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

//====================== Helper ================
func showImg(c *TopCameraAcquisition) {
	// now we can run stuff on the main thread like this
	mainthread.CallNonBlock(func() {
		prevWindow := gocv.NewWindow("Preview Window")
		defer prevWindow.Close()

		for {
			prevWindow.IMShow(c.Read())

			keyPressed := prevWindow.WaitKey(1)
			if keyPressed == 'q' {
				return
			}
		}
	})
}
