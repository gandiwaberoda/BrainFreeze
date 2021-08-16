package acquisition

import (
	"fmt"
	"image"
	"strconv"
	"sync"
	"time"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

type ForwardCameraAcquisition struct {
	Lock         *sync.RWMutex
	vc           *gocv.VideoCapture
	conf         *configuration.FreezeConfig
	postFrame    gocv.Mat
	postHsvFrame gocv.Mat
	firstFrame   bool
}

func NewForwardCameraAcquisition(conf *configuration.FreezeConfig) *ForwardCameraAcquisition {
	_postframe := gocv.NewMat()
	_postHsvFrame := gocv.NewMat()

	return &ForwardCameraAcquisition{
		conf:         conf,
		Lock:         &sync.RWMutex{},
		firstFrame:   false,
		postFrame:    _postframe,
		postHsvFrame: _postHsvFrame,
		vc:           &gocv.VideoCapture{},
	}
}

func workerForward(c *ForwardCameraAcquisition) {
	for {
		// TODO: Bisa dipindah lockingnya ke dalam read() function
		c.Lock.Lock()
		c.read()
		c.Lock.Unlock()
	}
}

func (c *ForwardCameraAcquisition) read() {
	// Baca frame dari kamera
	frame := gocv.NewMat()
	defer frame.Close()
	if ok := c.vc.Read(&frame); !ok {
		// Reopen vc
		c.vc.Close()
		c.Start()
		fmt.Println("Reopening forward camera because no next frame")
		return
	}

	// if frame.Empty() {
	// 	return
	// }

	// // Masking area lingkaran
	// circleMask := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)
	// defer circleMask.Close()
	// mid := image.Point{
	// 	X: c.conf.Camera.MidpointX, Y: c.conf.Camera.MidpointY,
	// }
	// white := color.RGBA{255, 255, 255, 0}
	// gocv.Circle(&circleMask, mid, c.conf.Camera.MidpointRad, white, -1)

	// maskedframe := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), gocv.MatTypeCV8U)
	// defer maskedframe.Close()
	// frame.CopyToWithMask(&maskedframe, circleMask)

	// // Ambil area persegi ROI
	// x0 := c.conf.Camera.MidpointX - c.conf.Camera.MidpointRad
	// y0 := c.conf.Camera.MidpointY - c.conf.Camera.MidpointRad
	// x1 := c.conf.Camera.MidpointX + c.conf.Camera.MidpointRad
	// y1 := c.conf.Camera.MidpointY + c.conf.Camera.MidpointRad
	// rect := image.Rect(x0, y0, x1, y1)
	// resImg := maskedframe.Region(rect)
	// defer resImg.Close()

	// // Flip vertically
	// gocv.Flip(resImg, &resImg, 0)

	// Normalize ukuran biar standar
	newSize := image.Point{c.conf.Camera.ForPostWidth, c.conf.Camera.ForPostHeight}
	gocv.Resize(frame, &c.postFrame, newSize, 0, 0, gocv.InterpolationLinear)

	gocv.CvtColor(c.postFrame, &c.postHsvFrame, gocv.ColorBGRToHSV)
	gocv.GaussianBlur(c.postHsvFrame, &c.postHsvFrame, image.Point{7, 7}, 0, 0, gocv.BorderDefault)

	c.firstFrame = true
}

func (c *ForwardCameraAcquisition) Read(dst *gocv.Mat) {
	if c.postFrame.Empty() {
		<-time.After(time.Millisecond * 1000)
		fmt.Println("Waiting forward camera...")
		c.Read(dst)
	} else {
		c.postFrame.CopyTo(dst)
	}
}

func (c *ForwardCameraAcquisition) ReadHSV(dst *gocv.Mat) {
	fmt.Println("col", c.postFrame.Cols())
	fmt.Println("row", c.postFrame.Rows())

	if !c.firstFrame {
		<-time.After(time.Millisecond * 1000)
	}

	if c.postHsvFrame.Empty() {
		fmt.Println("Waiting forward camera hsv...")
		c.ReadHSV(dst)
	} else {
		c.postHsvFrame.CopyTo(dst)
	}
}

func (c *ForwardCameraAcquisition) Start() {
	c.firstFrame = false

	src := c.conf.Camera.SrcForward[0]

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
		c.vc, errVc = gocv.VideoCaptureFile(c.conf.Camera.SrcForward[0])
	}
	if errVc != nil {
		panic(errVc)
	}

	go workerForward(c)
}

func (c *ForwardCameraAcquisition) Stop() {
	c.vc.Close()
}
