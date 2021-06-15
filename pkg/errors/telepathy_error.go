package frerror

import "fmt"

type TelepathyError struct {
	Detail string
}

func (e *TelepathyError) Error() string {
	return fmt.Sprint("General telepathy error", e.Detail)
}

type TelepathyNotRunningError struct {
	Detail string
}

func (e *TelepathyNotRunningError) Error() string {
	return fmt.Sprint("Telepathy is not running", e.Detail)
}
