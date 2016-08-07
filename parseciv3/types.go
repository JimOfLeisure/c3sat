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

// Pass it a ListItem struct so it can properly build an array and parse
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
	// list.List = make([]ListItem, list.Count)
	// for i := int32(0); i < list.Count; i++ {
	// 	err = binary.Read(r, binary.LittleEndian, &e)
	// 	list.List[i] = e
	// 	if err != nil {
	// 		return list, ReadError{err}
	// 	}
	// }

	// var temp BicResources
	// // err = binary.Read(r, binary.LittleEndian, &temp)
	// err = binary.Read(r, binary.LittleEndian, e)
	// if err != nil {
	// 	return list, ReadError{err}
	// }
	// fmt.Printf("temp %v", temp)

	// list.List = make([]ListItem, list.Count)
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
	// err = binary.Read(r, binary.LittleEndian, temp)
	return list, nil
}

// BicResources is part of the second SAV file section. Guessing at the alignment
type BicResources struct {
	A            int32
	ResourcePath [0x100]byte
	B            int32
	BicPath      [0x100]byte
	C            int32
}
