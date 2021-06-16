package acquisition

import (
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock       *sync.RWMutex
	vc         *gocv.VideoCapture
	conf       *configuration.FreezeConfig
	frame      gocv.Mat
	firstFrame bool
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig, threaded bool, preprocess bool, previewImage bool) *TopCameraAcquisition {
	return &TopCameraAcquisition{
		conf:       conf,
		Lock:       &sync.RWMutex{},
		firstFrame: false,
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

	frameCropped := res.Region(rect)

	newSize := image.Point{c.conf.Camera.PostWidth, c.conf.Camera.PostHeight}
	gocv.Resize(frameCropped, &res, newSize, 0, 0, gocv.InterpolationLinear)
	c.frame = res
	c.firstFrame = true
}

func (c *TopCameraAcquisition) Read() gocv.Mat {
	if !c.firstFrame {
		<-time.After(time.Millisecond * 200)
	}
	return c.frame
}

func (c *TopCameraAcquisition) Start() {
	c.firstFrame = false

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
}

func (c *TopCameraAcquisition) Stop() {
	c.vc.Close()
}
