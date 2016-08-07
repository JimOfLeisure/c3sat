package parseciv3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const debugContextBytes int = 0x300

// no err return because I'm calling this from inside errors
func debugHexDump(r *bytes.Reader) string {
	s := make([]byte, debugContextBytes)
	_ = binary.Read(r, binary.LittleEndian, &s)
	return hex.Dump(s)
}

// FileError returns errors while trying to open or decompress the file. Pass it the downstream error e.g. return FileError{err}
type FileError struct {
	e error
}

func (e FileError) Error() string {
	return fmt.Sprintf("Error opening or decompressing file: %s", e.e.Error())
}

// ReadError returns errors while trying to read data for parsing. Pass it the downstream error e.g. return ReadError{err}
type ReadError struct {
	e error
}

func (e ReadError) Error() string {
	return fmt.Sprintf("Error reading data: %s", e.e.Error())
}

// ParseError is when the data does not match an expected pattern. Pass it message string, expected value and hex dump of pertinent data.
type ParseError struct {
	s, Expected, Hexdump string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("Error parsing: %s", e.s)
}
