package civ3decompress

import (
	"fmt"
)

// FileError returns errors while trying to open or decompress the file. Pass it the downstream error e.g. return FileError{err}
type FileError struct {
	e error
}

func (e FileError) Error() string {
	return fmt.Sprintf("Error reading file: %s", e.e.Error())
}

// DecodeError is when the data does not match an expected pattern. Pass it message string.
type DecodeError struct {
	s string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("Error decoding: %s", e.s)
}
