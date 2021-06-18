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
	Lock        *sync.RWMutex
	vc          *gocv.VideoCapture
	conf        *configuration.FreezeConfig
	frame       gocv.Mat
	postFrame   gocv.Mat
	firstFrame  bool
	maskedFrame gocv.Mat
	circleMask  gocv.Mat
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig) *TopCameraAcquisition {
	circleMask := gocv.NewMatWithSize(conf.Camera.RawHeight, conf.Camera.RawWidth, gocv.MatTypeCV8UC3)

	mid := image.Point{
		X: conf.Camera.MidpointX, Y: conf.Camera.MidpointY,
	}
	white := color.RGBA{255, 255, 255, 0}
	gocv.Circle(&circleMask, mid, conf.Camera.MidpointRad, white, -1)

	_rawframe := gocv.NewMatWithSize(conf.Camera.RawHeight, conf.Camera.RawWidth, gocv.MatTypeCV8UC3)
	_maskedframe := gocv.NewMatWithSize(conf.Camera.RawHeight, conf.Camera.RawWidth, gocv.MatTypeCV8UC3)
	_postframe := gocv.NewMatWithSize(conf.Camera.PostHeight, conf.Camera.PostWidth, gocv.MatTypeCV8UC3)

	return &TopCameraAcquisition{
		conf:        conf,
		Lock:        &sync.RWMutex{},
		firstFrame:  false,
		frame:       _rawframe,
		maskedFrame: _maskedframe,
		postFrame:   _postframe,
		circleMask:  circleMask,
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
	// Baca frame dari kamera
	c.vc.Read(&c.frame)

	// Masking area lingkaran
	gocv.BitwiseAnd(c.frame, c.circleMask, &c.maskedFrame)

	x0 := c.conf.Camera.MidpointX - c.conf.Camera.MidpointRad
	y0 := c.conf.Camera.MidpointY - c.conf.Camera.MidpointRad
	x1 := c.conf.Camera.MidpointX + c.conf.Camera.MidpointRad
	y1 := c.conf.Camera.MidpointY + c.conf.Camera.MidpointRad
	rect := image.Rect(x0, y0, x1, y1)

	// Ambil area persegi ROI
	resImg := c.maskedFrame.Region(rect)
	defer resImg.Close()

	cropped := resImg.Clone()
	defer cropped.Close()

	// Normalize ukuran biar standar di hsv sama dnn
	newSize := image.Point{c.conf.Camera.PostWidth, c.conf.Camera.PostHeight}
	gocv.Resize(cropped, &c.postFrame, newSize, 0, 0, gocv.InterpolationLinear)
	c.firstFrame = true
}

func (c *TopCameraAcquisition) Read(dst *gocv.Mat) {
	if !c.firstFrame {
		<-time.After(time.Millisecond * 1000)
	}

	c.postFrame.CopyTo(dst)
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
