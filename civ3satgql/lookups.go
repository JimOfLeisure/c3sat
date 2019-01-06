package civ3satgql

import (
	"errors"
	"strconv"
)

// to make calling functions readable
const signed = true
const unsigned = false

func readInt32(offset int, signed bool) int {
	n := int(saveGame.data[offset]) +
		int(saveGame.data[offset+1])*0x100 +
		int(saveGame.data[offset+2])*0x10000 +
		int(saveGame.data[offset+3])*0x1000000
	if signed && n > 0x80000000 {
		n = n - 0x80000000
	}
	return n
}

func readInt16(offset int, signed bool) int {
	n := int(saveGame.data[offset]) +
		int(saveGame.data[offset+1])*0x100
	if signed && n > 0x8000 {
		n = n - 0x8000
	}
	return n
}

func readInt8(offset int, signed bool) int {
	n := int(saveGame.data[offset])
	if signed && n > 0x80 {
		n = n - 0x80
	}
	return n
}

func sectionOffset(sectionName string, nth int) (int, error) {
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
