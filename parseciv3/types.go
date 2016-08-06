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
	parse(r io.ReadSeeker) error
}

// Civ3 is the SAV file header
type Civ3 struct {
	Name [4]byte
	// 28 bytes. Guessing on alignment
	A, B, C, D, E, F uint32
	G                uint16
}

func (me *Civ3) parse(r io.ReadSeeker) error {
	// var temp Civ3
	// _, err := binary.ReadSeeker
	err := binary.Read(r, binary.LittleEndian, me)
	if err != nil {
		return ReadError{err}
	}
	return nil
}
