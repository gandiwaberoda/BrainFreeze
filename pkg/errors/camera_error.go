package freeze_error

import "errors"

func CameraError() error {
	return errors.New("Camera failed to start")
}
