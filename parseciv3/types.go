package parseciv3

import (
	"encoding/binary"
	"io"
)

// ParsedData is the structure of the parsed data
type ParsedData map[string]Section

// Civ3Data contains the game data
type Civ3Data struct {
	FileName   string
	Compressed bool
	Data       ParsedData
	Next       string
	// RawFile    []byte
}

// Section is the inteface for the various structs decoded from the data files
type Section interface {
	Name() string
}

// Civ3 is the SAV file header
type Civ3 struct {
	ClassName [4]byte
	// 28 bytes. Guessing on alignment
	A, B, C, D, E, F uint32
	G                uint16
}

func (s Civ3) Name() string {
	return string(s.ClassName[:])
}

func newCiv3(r io.ReadSeeker) (Civ3, error) {
	var data Civ3
	err := binary.Read(r, binary.LittleEndian, &data)
	if err != nil {
		return data, ReadError{err}
	}
	return data, nil
}
