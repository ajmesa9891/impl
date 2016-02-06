package impl

import (
	"fmt"
)

type InvalidInterfacePathError struct {
	message string
}

func NewInvalidInterfacePathError(message string, args ...interface{}) *InvalidInterfacePathError {
	return &InvalidInterfacePathError{fmt.Sprintf(message, args...)}
}

func (e *InvalidInterfacePathError) Error() string {
	return e.message
}
