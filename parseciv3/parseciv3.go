package parseciv3

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	civ3data.FileName = path

	// Load file into struct for parsing
	rawFile, compressed, err := ReadFile(path)
	if err != nil {
		return civ3data, err
	}
	civ3data.Compressed = compressed

	// Create ReadSeeker
	r := bytes.NewReader(rawFile)

	civ3data.Data, err = ParseCiv3(r)
	if err != nil {
		return civ3data, err
	}
	civ3data.Next = debugHexDump(r)

	// // TEMP writing 2nd GAME portion out to file for analysis
	// // from http://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	// outFileName := path + ".game"
	// if _, err := os.Stat(outFileName); os.IsNotExist(err) {
	// 	mybytes := make([]byte, 0x1200)
	// 	_, _ = r.Read(mybytes)
	// 	_ = ioutil.WriteFile(outFileName, mybytes, 0644)
	// }
	return civ3data, nil
}

// peek returns the next 4 bytes nondestructively
func peek(r io.ReadSeeker) ([]byte, error) {
	var b [4]byte
	err := binary.Read(r, binary.LittleEndian, &b)
	if err != nil {
		return nil, ReadError{err}
	}
	// Back the pointer up 4 bytes
	r.Seek(-4, 1)
	return b[:], nil
}

// peekName is a wrapper to allow 4-char strings as a parameter to peekFour
func peekName(r io.ReadSeeker, expected string) error {
	var b [4]byte
	copy(b[:], expected)
	return peekFour(r, b)
}

// peekFour returns r with the pointer in its original location, and an error if the next bytes don't match expected
func peekFour(r io.ReadSeeker, expected [4]byte) error {
	var peek [4]byte
	err := binary.Read(r, binary.LittleEndian, &peek)
	if err != nil {
		return ReadError{err}
	}
	// Back the pointer up 4 bytes
	r.Seek(-4, 1)
	if peek != expected {
		return ParseError{"Parse error: Unexpected data", fmt.Sprintf("%v", expected), debugHexDump(r)}
	}
	return nil
}

// ParseCiv3 takes raw save file data and returns a map of the parsed data
func ParseCiv3(r io.ReadSeeker) (ParsedData, error) {
	data := make(ParsedData)
	var err error

	// CIV3 section, optional if BIC file
	err = peekName(r, "CIV3")
	if err == nil {
		// need to preserve parent scope's err
		var err error
		data["CIV3"], err = newCiv3(r)
		if err != nil {
			return data, err
		}
		// BIC_
		data["BIC "], err = newBase(r)
		if err != nil {
			return data, err
		}
	} else {
		if _, ok := err.(ParseError); ok {
			// continue if not matched; may be a BIC/X/Q file
		} else {
			return data, err
		}
	}

	// BIC file / section
	bicHeader := make([]byte, 4)
	err = binary.Read(r, binary.LittleEndian, &bicHeader)
	if err != nil {
		return data, err
	}
	switch string(bicHeader) {
	case "BIC ", "BICX", "BICQ":
		// continue
	default:
		// Back the pointer up 4 bytes
		r.Seek(-4, 1)
		return data, ParseError{"Parse error: Unexpected data", "BIC*", debugHexDump(r)}

	}
	var gameSectionCount int
	// TODO: Add sections for custom world map
	// loop sections until GAME reached
	// for name, err := peek(r); string(name[:]) != "GAME"; name, err = peek(r) {
	for name, err := peek(r); gameSectionCount < 1 || string(name[:]) != "GAME"; name, err = peek(r) {
		if err != nil {
			return data, err
		}
		switch string(name[:]) {
		case "GAME":
			switch gameSectionCount {
			case 0:
				data["GAME"], err = newList(r)
				// case 1:
				// 	// this is not right
				// 	data["GAME2"], err = newBase(r)
			}
			if err != nil {
				return data, err
			}
			gameSectionCount++
			// fallthrough
		// (Almost?) Always in this order, but have seen FLAV after EXPR in some saves and after WSIZ in others
		case "VER#", "BLDG", "CTZN", "CULT", "DIFF", "ERAS", "ESPN", "EXPR", "GOOD", "GOVT", "RULE", "PRTO", "RACE", "TECH", "TFRM", "TERR", "WSIZ", "LEAD":
			data[string(name[:])], err = newList(r)
			if err != nil {
				return data, err
			}
		case "FLAV":
			data[string(name[:])], err = newFlav(r)
			if err != nil {
				return data, err
			}
		default:
			return data, ParseError{"Parse error: Unexpected data", "<known classname>", debugHexDump(r)}
		}
	}
	var gameSection Game
	// gameSection := Game{}
	err = binary.Read(r, binary.LittleEndian, &gameSection)
	if err != nil {
		return data, err
	}
	data["GAME2"] = gameSection

	searchLength := 0x1200
	buffer := make([]byte, searchLength)
	_, err = r.Read(buffer)
	if err != nil {
		return data, err
	}
	index := bytes.Index(buffer, []byte("DATE"))
	fmt.Println(index)
	_, err = r.Seek(int64(index-searchLength), 1)
	if err != nil {
		return data, err
	}
	// var gameNext GameNext
	// err = binary.Read(r, binary.LittleEndian, &gameNext)
	// data["GameNext"] = gameNext
	return data, nil
}
