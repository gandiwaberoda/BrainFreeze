package frerror

import "fmt"

type SerialError struct {
	Detail string
}

func (e *SerialError) Error() string {
	return fmt.Sprint("General serial error", e.Detail)
}
