package errs

import (
	"fmt"
)

type EmptyInterfacePathError struct {
	message string
}

func NewEmptyInterfacePathError(message string, args ...interface{}) *EmptyInterfacePathError {
	return &EmptyInterfacePathError{fmt.Sprintf(message, args...)}
}

func (e *EmptyInterfacePathError) Error() string {
	return e.message
}

type InvalidInterfacePathError struct {
	message string
}

func NewInvalidInterfacePathError(message string, args ...interface{}) *InvalidInterfacePathError {
	return &InvalidInterfacePathError{fmt.Sprintf(message, args...)}
}

func (e *InvalidInterfacePathError) Error() string {
	return e.message
}

type InvalidImportFormatError struct {
	message string
}

func NewInvalidImportFormatError(message string, args ...interface{}) *InvalidImportFormatError {
	return &InvalidImportFormatError{fmt.Sprintf(message, args...)}
}

func (e *InvalidImportFormatError) Error() string {
	return e.message
}

type CouldNotFindPackageError struct {
	message string
}

func NewCouldNotFindPackageError(message string, args ...interface{}) *CouldNotFindPackageError {
	return &CouldNotFindPackageError{fmt.Sprintf(message, args...)}
}

func (e *CouldNotFindPackageError) Error() string {
	return e.message
}
