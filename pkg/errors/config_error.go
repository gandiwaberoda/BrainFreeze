package frerror

import "fmt"

type ConfigError struct {
	Detail string
}

func (e *ConfigError) Error() string {
	return fmt.Sprint("General local config error", e.Detail)
}
