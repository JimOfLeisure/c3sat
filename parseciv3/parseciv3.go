package parseciv3

// Copyright (c) 2016 Jim Nelson

// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:

// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"

	"github.com/myjimnelson/c3sat/civ3decompress"
)

// ReadFile takes a filename and returns the decompressed file data or the raw data if it's not compressed. Also returns true if compressed.
func ReadFile(path string) ([]byte, bool, error) {
	// Open file, hanlde errors, defer close
	file, err := os.Open(path)
	if err != nil {
		return nil, false, FileError{err}
	}
	defer file.Close()

	var compressed bool
	var data []byte
	header := make([]byte, 2)
	_, err = file.Read(header)
	if err != nil {
		return nil, false, FileError{err}
	}
	// reset pointer to parse from beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, false, FileError{err}
	}
	switch {
	case header[0] == 0x00 && (header[1] == 0x04 || header[1] == 0x05 || header[1] == 0x06):
		compressed = true
		data, err = civ3decompress.Decompress(file)
		if err != nil {
			return nil, false, err
		}
	default:
		// log.Println("Not a compressed file. Proceeding with uncompressed stream.")
		// TODO: I'm sure I'm doing this in a terribly inefficient way. Need to refactor everything to pass around file pointers I think
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, false, FileError{err}
		}
	}
	return data, compressed, error(nil)

}

// TODO: Do I really need or want rawFile in the struct?
// A: yes, while decoding, anyway. Or maybe I can include a lookahead dumping field instead of the entire file?

// NOTE: Just changed ParseCiv3 to NewCiv3Data. It sets up the struct, and then I'll call a new ParseSav and/or ParseBic to populate the map field

// NewCiv3Data takes a path to a file and returns a struct containing the parsed data and a rawFile field
func NewCiv3Data(path string) (Civ3Data, error) {
	var civ3data Civ3Data
	var compressed bool
	var err error
	var listHeader ListHeader
	civ3data.FileName = path

	// Load file into struct for parsing
	rawFile, compressed, err := ReadFile(path)
	if err != nil {
		return civ3data, err
	}
	civ3data.Compressed = compressed

	// Create ReadSeeker
	r := bytes.NewReader(rawFile)

	// CIV3 section, optional if BIC file
	name, _, err := peek(r)
<<<<<<< HEAD
	if err != nil {
		return civ3data, err
	}
	if string(name[:]) == "CIV3" {
		err = binary.Read(r, binary.LittleEndian, &civ3data.Civ3)
		if err != nil {
			return civ3data, err
		}

		// BIC_ resoruces
		err = binary.Read(r, binary.LittleEndian, &civ3data.BicResources)
		if err != nil {
			return civ3data, err
		}
	}
	// BIC file / section
	err = binary.Read(r, binary.LittleEndian, &civ3data.BicFileHeader)
	if err != nil {
		return civ3data, err
	}
	switch string(civ3data.BicFileHeader[:]) {
	case "BIC ", "BICX", "BICQ":
		// continue
	default:
		return civ3data, ParseError{"Parse error: Unexpected data", "BIC*", debugHexDump(r)}

	}
	// VER#
	err = binary.Read(r, binary.LittleEndian, &listHeader)
	if err != nil {
		return civ3data, err
	}
	civ3data.VerNum = make([]VerNum, listHeader.Count)
	err = binary.Read(r, binary.LittleEndian, &civ3data.VerNum)
	if err != nil {
		return civ3data, err
	}

	// Custom rules
	err = binary.Read(r, binary.LittleEndian, &listHeader)
	if err != nil {
		return civ3data, err
	}
	if string(listHeader.Name[:]) == "BLDG" {
		civ3data.Bldg = make([]Bldg, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Bldg)
		if err != nil {
			return civ3data, err
		}
	}
	// Original parseddata section
	civ3data.Data, err = ParseCiv3(r)
=======
>>>>>>> jim
	if err != nil {
		return civ3data, err
	}
	if string(name[:]) == "CIV3" {
		err = binary.Read(r, binary.LittleEndian, &civ3data.Civ3)
		if err != nil {
			return civ3data, err
		}

<<<<<<< HEAD
// peek returns the next 4 bytes nondestructively
func peek(r io.ReadSeeker) ([]byte, int32, error) {
	var b [4]byte
	err := binary.Read(r, binary.LittleEndian, &b)
	if err != nil {
		return nil, 0, ReadError{err}
	}
	var length int32
	err = binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return b[:], 0, ReadError{err}
	}
	// Back the pointer up 4 bytes
	r.Seek(-8, 1)
	return b[:], length, nil
}

=======
		// BIC_ resoruces
		err = binary.Read(r, binary.LittleEndian, &civ3data.BicResources)
		if err != nil {
			return civ3data, err
		}
	}
	// BIC file / section
	err = binary.Read(r, binary.LittleEndian, &civ3data.BicFileHeader)
	if err != nil {
		return civ3data, err
	}
	switch string(civ3data.BicFileHeader[:]) {
	case "BIC ", "BICX", "BICQ":
		// continue
	default:
		return civ3data, ParseError{"Parse error: Unexpected data", "BIC*", debugHexDump(r)}

	}
	// VER#
	err = binary.Read(r, binary.LittleEndian, &listHeader)
	if err != nil {
		return civ3data, err
	}
	civ3data.VerNum = make([]VerNum, listHeader.Count)
	err = binary.Read(r, binary.LittleEndian, &civ3data.VerNum)
	if err != nil {
		return civ3data, err
	}

	// Custom rules
	err = binary.Read(r, binary.LittleEndian, &listHeader)
	if err != nil {
		return civ3data, err
	}
	if string(listHeader.Name[:]) == "BLDG" {
		// BLDG
		civ3data.Bldg = make([]Bldg, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Bldg)
		if err != nil {
			return civ3data, err
		}
		// CTZN
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Ctzn = make([]Ctzn, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Ctzn)
		if err != nil {
			return civ3data, err
		}
		// CULT
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Cult = make([]Cult, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Cult)
		if err != nil {
			return civ3data, err
		}
		// DIFF
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Diff = make([]Difficulty, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Diff)
		if err != nil {
			return civ3data, err
		}
		// ERAS
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Eras = make([]Era, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Eras)
		if err != nil {
			return civ3data, err
		}
		// ESPN
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Espn = make([]Espn, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Espn)
		if err != nil {
			return civ3data, err
		}
		// EXPR
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Expr = make([]Expr, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Expr)
		if err != nil {
			return civ3data, err
		}
		// FLAV
		// reading this one in incrementally because it's an array of arrays, and the leaves have an array
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		civ3data.Flav = make([][]Flavor, listHeader.Count)
		for i := range civ3data.Flav {
			var numFlavors int32
			err = binary.Read(r, binary.LittleEndian, &numFlavors)
			if err != nil {
				return civ3data, err
			}
			civ3data.Flav[i] = make([]Flavor, numFlavors)
			for j := range civ3data.Flav[i] {
				err = binary.Read(r, binary.LittleEndian, &civ3data.Flav[i][j])
				if err != nil {
					return civ3data, err
				}
			}
		}
		// GOOD
		err = binary.Read(r, binary.LittleEndian, &listHeader)
		if err != nil {
			return civ3data, err
		}
		// fmt.Println(string(listHeader.Name[:]))
		civ3data.Good = make([]Good, listHeader.Count)
		err = binary.Read(r, binary.LittleEndian, &civ3data.Good)
		if err != nil {
			return civ3data, err
		}
	} else {
		// No custom rules, so back up pointer before going to old parseddata function
		_, err = r.Seek(-8, 1)
		if err != nil {
			return civ3data, err
		}
	}
	// Original parseddata section
	civ3data.Data, err = ParseCiv3(r)
	if err != nil {
		return civ3data, err
	}

	err = binary.Read(r, binary.LittleEndian, &civ3data.Wrld)
	if err != nil {
		return civ3data, err
	}
	tileCount := civ3data.Wrld.MapHeight * int32(civ3data.Wrld.MapWidth/2)
	_ = tileCount
	civ3data.Tile = make([]Tile, tileCount)
	err = binary.Read(r, binary.LittleEndian, &civ3data.Tile)
	if err != nil {
		return civ3data, err
	}
	civ3data.Cont = make([]Continent, civ3data.Wrld.NumContinents)
	err = binary.Read(r, binary.LittleEndian, &civ3data.Cont)
	if err != nil {
		return civ3data, err
	}

	// TODO: Find where resource count is
	// Currently counting GOODs if present or using default 26 resource count
	var numResources int
	if len(civ3data.Good) > 0 {
		numResources = len(civ3data.Good)
	} else {
		numResources = 26
	}
	civ3data.ResourceCounts = make([]int32, numResources)
	err = binary.Read(r, binary.LittleEndian, &civ3data.ResourceCounts)
	if err != nil {
		return civ3data, err
	}
	civ3data.Next = debugHexDump(r)
	return civ3data, nil
}

// peek returns the next 4 bytes nondestructively
func peek(r io.ReadSeeker) ([]byte, int32, error) {
	var b [4]byte
	err := binary.Read(r, binary.LittleEndian, &b)
	if err != nil {
		return nil, 0, ReadError{err}
	}
	var length int32
	err = binary.Read(r, binary.LittleEndian, &length)
	if err != nil {
		return b[:], 0, ReadError{err}
	}
	// Back the pointer up 4 bytes
	r.Seek(-8, 1)
	return b[:], length, nil
}

>>>>>>> jim
// ParseCiv3 takes raw save file data and returns a map of the parsed data
func ParseCiv3(r io.ReadSeeker) (ParsedData, error) {
	data := make(ParsedData)
	var gameSectionCount int
	// TODO: Add sections for custom world map
	// loop sections until 2nd GAME reached which marks end of the BIC and beginning of game data
	for name, length, err := peek(r); gameSectionCount < 1 || string(name[:]) != "GAME"; name, length, err = peek(r) {
		_ = length
		if err != nil {
			return data, err
		}
		switch string(name[:]) {
		case "GAME":
			data["GAME"], err = newList(r)
			if err != nil {
				return data, err
			}
			gameSectionCount++
		// (Almost?) Always in this order, but have seen FLAV after EXPR in some saves and after WSIZ in others
		case "GOVT", "RULE", "PRTO", "RACE", "TECH", "TFRM", "TERR", "WSIZ", "LEAD":
			data[string(name[:])], err = newList(r)
			if err != nil {
				return data, err
			}
		default:
			return data, ParseError{"Parse error: Unexpected data", "<known classname>", debugHexDump(r)}
		}
	}
	var abort bool
	// 2nd GAME / first non-BIC GAME and beyond
	for byteName, length, err := peek(r); !abort; byteName, length, err = peek(r) {
		_ = length
		if err != nil {
			return data, err
		}
		name := string(byteName)
		switch name {
		case "GAME":
			var gameSection Game
			err = binary.Read(r, binary.LittleEndian, &gameSection)
			if err != nil {
				return data, err
			}
			data["GAME2"] = gameSection

			searchLength := 0x4000
			buffer := make([]byte, searchLength)
			_, err = r.Read(buffer)
			if err != nil {
				return data, err
			}
			index := bytes.Index(buffer, []byte("DATE"))
			_, err = r.Seek(int64(-searchLength), 1)
			if err != nil {
				return data, err
			}
			buffer = make([]byte, index)
			_, err = r.Read(buffer)
			if err != nil {
				return data, err
			}
			data["WTF"] = hex.Dump(buffer)

		case "DATE":
			// briefly fetching the DATEs and PLGIs this way .. and two ints between that and CNSL
			for i := 1; i < 6; i++ {
				data[string(i)], err = newBase(r)
				if err != nil {
					return data, err
				}
			}
			intBuffer := make([]int32, 2)
			err = binary.Read(r, binary.LittleEndian, &intBuffer)
			if err != nil {
				return data, err
			}
			data["BeforeCNSL"] = intBuffer
		case "CNSL":
			data["CNSL"], err = newBase(r)
			if err != nil {
				return data, err
			}
		default:
			abort = true
		}
	}
	return data, nil
}
