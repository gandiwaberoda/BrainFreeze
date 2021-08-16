package diagnostic

import (
	"errors"
	"fmt"
	"strconv"

	"gocv.io/x/gocv"
	"harianugrah.com/brainfreeze/pkg/models/configuration"
)

func ConfigValidate(c configuration.FreezeConfig) error {
	fmt.Println("Starting self diagnostic")

	err := _checkCameraRawConfig(c)
	if err != nil {
		return err
	}

	return nil
}

func _checkCameraRawConfig(c configuration.FreezeConfig) error {
	// Lingkaran mask tidak boleh melewati ukuran frame
	if c.Camera.MidpointX+c.Camera.MidpointRad > c.Camera.RawWidth {
		return errors.New("MidpointX+MidpointRad tidak boleh keluar RawWidth")
	}

	if c.Camera.MidpointX-c.Camera.MidpointRad < 0 {
		return errors.New("MidpointX-MidpointRad tidak boleh keluar 0")
	}

	if c.Camera.MidpointY-c.Camera.MidpointRad < 0 {
		return errors.New("MidpointX+MidpointRad tidak boleh keluar 0")
	}

	if c.Camera.MidpointY+c.Camera.MidpointRad > c.Camera.RawHeight {
		return errors.New("MidpointX+MidpointRad tidak boleh keluar RawHeight")
	}

	// TOP CAMERA TEST
	if err := CheckTopCamera(c); err != nil {
		return err
	}

	// FORWARD CAMERA TEST
	if err := CheckForwardCamera(c); err != nil {
		return err
	}

	return nil
}

func CheckTopCamera(c configuration.FreezeConfig) error {

	// Cek apakah RawHeight dan RawWidth sesuai dengan output dari VideoCapture
	src := c.Camera.Src[0]
	var errVc error
	var vc *gocv.VideoCapture

	if len(src) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(src)
		if errInt != nil {
			panic(errInt)
		}
		vc, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		vc, errVc = gocv.VideoCaptureFile(c.Camera.Src[0])
	}
	defer vc.Close()

	if errVc != nil {
		return errVc
	}

	firstFrame := gocv.NewMat()
	vc.Read(&firstFrame)
	defer firstFrame.Close()

	if firstFrame.Empty() {
		return errors.New("can't get first frame")
	}

	if firstFrame.Rows() != c.Camera.RawHeight {
		return errors.New("rawHeight tidak sama dengan output frame Rows " + strconv.Itoa(firstFrame.Rows()))
	}

	if firstFrame.Cols() != c.Camera.RawWidth {
		return errors.New("rawWidth tidak sama dengan output frame Cols " + strconv.Itoa(firstFrame.Cols()))
	}

	return nil
}

func CheckForwardCamera(c configuration.FreezeConfig) error {

	// Cek apakah RawHeight dan RawWidth sesuai dengan output dari VideoCapture
	src := c.Camera.SrcForward[0]
	var errVc error
	var vc *gocv.VideoCapture

	if len(src) == 1 {
		// Kamera
		srcInt, errInt := strconv.Atoi(src)
		if errInt != nil {
			panic(errInt)
		}
		vc, errVc = gocv.VideoCaptureDevice(srcInt)
	} else {
		// Video
		vc, errVc = gocv.VideoCaptureFile(src)
	}
	defer vc.Close()

	if errVc != nil {
		return errVc
	}

	firstFrame := gocv.NewMat()
	vc.Read(&firstFrame)
	defer firstFrame.Close()

	if firstFrame.Empty() {
		return errors.New("can't get forward camera first frame")
	}

	if firstFrame.Rows() != c.Camera.ForHeight {
		return errors.New("height forward camera tidak sama dengan output frame Rows " + strconv.Itoa(firstFrame.Rows()))
	}

	if firstFrame.Cols() != c.Camera.ForWidth {
		return errors.New("width tidak sama dengan output frame Cols " + strconv.Itoa(firstFrame.Cols()))
	}

	return nil
}
