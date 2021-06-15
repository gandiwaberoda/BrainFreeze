package frerror

import "fmt"

type CameraError struct {
	Detail string
}

func (e *CameraError) Error() string {
	return fmt.Sprint("General camera error", e.Detail)
}
