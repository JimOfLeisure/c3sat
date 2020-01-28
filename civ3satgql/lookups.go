package civ3satgql

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"golang.org/x/text/encoding/charmap"
)

// to make calling functions readable
const Signed = true
const Unsigned = false

//deprecating in favor of *saveGameType.readInt32()
func ReadInt32(offset int, signed bool) int {
	n := int(saveGame.data[offset]) +
		int(saveGame.data[offset+1])*0x100 +
		int(saveGame.data[offset+2])*0x10000 +
		int(saveGame.data[offset+3])*0x1000000
	if signed && n&0x80000000 != 0 {
		n = -(n ^ 0xffffffff + 1)
	}
	return n
}

//deprecating in favor of *saveGameType.readInt16()
func ReadInt16(offset int, signed bool) int {
	n := int(saveGame.data[offset]) +
		int(saveGame.data[offset+1])*0x100
	if signed && n&0x8000 != 0 {
		n = -(n ^ 0xffff + 1)
	}
	return n
}

//deprecating in favor of *saveGameType.readInt8()
func ReadInt8(offset int, signed bool) int {
	n := int(saveGame.data[offset])
	if signed && n&0x80 != 0 {
		n = -(n ^ 0xff + 1)
	}
	return n
}

//deprecating in favor of *saveGameType.sectionOffset()
func SectionOffset(sectionName string, nth int) (int, error) {
	var i, n int
	for i < len(saveGame.sections) {
		if saveGame.sections[i].name == sectionName {
			n++
			if n >= nth {
				return saveGame.sections[i].offset + len(sectionName), nil
			}
		}
		i++
	}
	return -1, errors.New("Could not find " + strconv.Itoa(nth) + " section named " + sectionName)
}

// CivString Finds null-terminated string and converts from Windows-1252 to UTF-8
func CivString(b []byte) (string, error) {
	var win1252 string
	for i := 0; i < len(b); i++ {
		if b[i] == 0 {
			win1252 = string(b[:i])
			break
		}
	}
	if win1252 == "" {
		win1252 = string(b)
	}
	sr := strings.NewReader(win1252)
	tr := charmap.Windows1252.NewDecoder().Reader(sr)

	outUtf8, err := ioutil.ReadAll(tr)
	if err != nil {
		return "", err
	}

	return string(outUtf8), nil
}
