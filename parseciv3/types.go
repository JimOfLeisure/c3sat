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
type Section interface{}

// ListItem are the structs in a list
type ListItem interface{}

// Civ3 is the SAV file header
type Civ3 struct {
	Name [4]byte
	// 28 bytes. Guessing on alignment
	A, B, C, D, E, F uint32
	G                uint16
}

func newCiv3(r io.ReadSeeker) (Civ3, error) {
	var data Civ3
	err := binary.Read(r, binary.LittleEndian, &data)
	if err != nil {
		return data, ReadError{err}
	}
	return data, nil
}

// Base is one of the basic section structures of the game data
type Base struct {
	Name    [4]byte
	Length  int32
	RawData []byte
}

func newBase(r io.ReadSeeker) (Base, error) {
	var base Base
	var err error
	err = binary.Read(r, binary.LittleEndian, &base.Name)
	if err != nil {
		return base, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &base.Length)
	if err != nil {
		return base, ReadError{err}
	}
	base.RawData = make([]byte, base.Length)
	err = binary.Read(r, binary.LittleEndian, &base.RawData)
	if err != nil {
		return base, ReadError{err}
	}
	return base, nil
}

// List is one of the basic section structures of the game data
type List struct {
	Name  [4]byte
	Count int32
	List  [][]byte
}

func newList(r io.ReadSeeker) (List, error) {
	var list List
	var err error
	err = binary.Read(r, binary.LittleEndian, &list.Name)
	if err != nil {
		return list, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &list.Count)
	if err != nil {
		return list, ReadError{err}
	}
	for i := int32(0); i < list.Count; i++ {
		var length int32
		err = binary.Read(r, binary.LittleEndian, &length)
		if err != nil {
			return list, ReadError{err}
		}

		temp := make([]byte, length)
		err = binary.Read(r, binary.LittleEndian, &temp)
		list.List = append(list.List, temp)

	}
	return list, nil
}

// Flav is one of the basic section structures of the game data
type Flav struct {
	Name  [4]byte
	Count int32
	List  [][]Flavor
}

func newFlav(r io.ReadSeeker) (Flav, error) {
	var flav Flav
	var err error
	err = binary.Read(r, binary.LittleEndian, &flav.Name)
	if err != nil {
		return flav, ReadError{err}
	}
	err = binary.Read(r, binary.LittleEndian, &flav.Count)
	if err != nil {
		return flav, ReadError{err}
	}
	for i := int32(0); i < flav.Count; i++ {
		var count int32
		err = binary.Read(r, binary.LittleEndian, &count)
		if err != nil {
			return flav, ReadError{err}
		}
		flavorGroups := make([]Flavor, count)
		flav.List = append(flav.List, flavorGroups)
		for j := int32(0); j < count; j++ {
			flav.List[i][j] = Flavor{}
			err = binary.Read(r, binary.LittleEndian, &flav.List[i][j])
			if err != nil {
				return flav, ReadError{err}
			}
		}
	}
	return flav, nil
}

// BicResources is part of the second SAV file section. Guessing at the alignment
type BicResources struct {
	A            int32
	ResourcePath [0x100]byte
	B            int32
	BicPath      [0x100]byte
	C            int32
}

type Flavor struct {
	A                      int32
	FlavorName             [0x100]byte
	B, C, D, E, F, G, H, I int32
}
