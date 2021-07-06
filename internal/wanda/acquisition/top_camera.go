package acquisition

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type TopCameraAcquisition struct {
	Lock          *sync.RWMutex
	vc            *gocv.VideoCapture
	conf          *configuration.FreezeConfig
	postFrame     gocv.Mat
	postHsvFrame  gocv.Mat
	postGrayFrame gocv.Mat
	firstFrame    bool
}

func CreateTopCameraAcquisition(conf *configuration.FreezeConfig) *TopCameraAcquisition {
	_postframe := gocv.NewMatWithSize(conf.Camera.PostHeight, conf.Camera.PostWidth, gocv.MatTypeCV8U)
	_postHsvFrame := gocv.NewMat()
	_postGrayFrame := gocv.NewMat()

	return &TopCameraAcquisition{
		conf:          conf,
		Lock:          &sync.RWMutex{},
		firstFrame:    false,
		postFrame:     _postframe,
		postHsvFrame:  _postHsvFrame,
		postGrayFrame: _postGrayFrame,
		vc:            &gocv.VideoCapture{},
	}
}

func workerTop(c *TopCameraAcquisition) {
	for {
		// TODO: Bisa dipindah lockingnya ke dalam read() function
		c.Lock.Lock()
		c.read()
		c.Lock.Unlock()
	}
}

func (c *TopCameraAcquisition) read() {
	// Baca frame dari kamera
	frame := gocv.NewMat()
	defer frame.Close()
	c.vc.Read(&frame)

	if ok := c.vc.Read(&frame); !ok {
		// Reopen vc
		c.vc.Close()
		c.Start()
		fmt.Println("Reopening top camera because no next frame")
		return
	}

	// Masking area lingkaran
	circleMask := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)
	defer circleMask.Close()
	mid := image.Point{
		X: c.conf.Camera.MidpointX, Y: c.conf.Camera.MidpointY,
	}
	white := color.RGBA{255, 255, 255, 0}
	gocv.Circle(&circleMask, mid, c.conf.Camera.MidpointRad, white, -1)

	maskedframe := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)
	defer maskedframe.Close()
	frame.CopyToWithMask(&maskedframe, circleMask)

	// Ambil area persegi ROI
	x0 := c.conf.Camera.MidpointX - c.conf.Camera.MidpointRad
	y0 := c.conf.Camera.MidpointY - c.conf.Camera.MidpointRad
	x1 := c.conf.Camera.MidpointX + c.conf.Camera.MidpointRad
	y1 := c.conf.Camera.MidpointY + c.conf.Camera.MidpointRad
	rect := image.Rect(x0, y0, x1, y1)
	resImg := maskedframe.Region(rect)
	defer resImg.Close()

	// Flip vertically
	gocv.Flip(resImg, &resImg, 0)

	// Normalize ukuran biar standar di hsv sama dnn
	newSize := image.Point{c.conf.Camera.PostWidth, c.conf.Camera.PostHeight}
	gocv.Resize(resImg, &c.postFrame, newSize, 0, 0, gocv.InterpolationLinear)

	// HSV Frame
	gocv.CvtColor(c.postFrame, &c.postHsvFrame, gocv.ColorBGRToHSV)
	gocv.GaussianBlur(c.postHsvFrame, &c.postHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

	// Gray Frame
	gocv.CvtColor(c.postFrame, &c.postGrayFrame, gocv.ColorBGRToGray)
	gocv.GaussianBlur(c.postGrayFrame, &c.postGrayFrame, image.Pt(3, 3), 0, 0, gocv.BorderDefault)
	gocv.EqualizeHist(c.postGrayFrame, &c.postGrayFrame)

	c.firstFrame = true
}

func (c *TopCameraAcquisition) Read(dst *gocv.Mat) {
	if !c.firstFrame {
		<-time.After(time.Millisecond * 1000)
	}

	if c.postFrame.Empty() {
		fmt.Println("Waiting top camera...")
		c.Read(dst)
	} else {
		c.postFrame.CopyTo(dst)
	}
}

func (c *TopCameraAcquisition) ReadHSV(dst *gocv.Mat) {
	if !c.firstFrame {
		<-time.After(time.Millisecond * 1000)
	}

	if c.postHsvFrame.Empty() {
		fmt.Println("Waiting top camera hsv...")
		c.ReadHSV(dst)
	} else {
		c.postHsvFrame.CopyTo(dst)
	}
}

func (c *TopCameraAcquisition) ReadGray(dst *gocv.Mat) {
	if !c.firstFrame {
		<-time.After(time.Millisecond * 1000)
	}

	if c.postHsvFrame.Empty() {
		fmt.Println("Waiting top camera gray...")
		c.ReadGray(dst)
	} else {
		c.postGrayFrame.CopyTo(dst)
	}
}

func (c *TopCameraAcquisition) Start() {
	c.firstFrame = false

	src := c.conf.Camera.Src[0]

	var errVc error
	if len(src) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(src)
		if errInt != nil {
			panic(errInt)
		}
		c.vc, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		c.vc, errVc = gocv.VideoCaptureFile(c.conf.Camera.Src[0])
	}
	if errVc != nil {
		panic(errVc)
	}

	go workerTop(c)
}

func (c *TopCameraAcquisition) Stop() {
	c.vc.Close()
}
