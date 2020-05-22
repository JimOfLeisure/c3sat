package luaciv3

import (
	"errors"
	"strconv"

	"github.com/myjimnelson/c3sat/civ3decompress"
)

// starting by copying wholesale parts of /queryciv3/queryciv3.go

// Signed - to make calling functions readable
const Signed = true

// Unsigned - to make calling functions readable
const Unsigned = false

type sectionType struct {
	name   string
	offset int
}

type saveGameType struct {
	path     string
	data     []byte
	sections []sectionType
}

var saveGame saveGameType
var defaultBic saveGameType
var currentBic saveGameType
var currentGame saveGameType

// populates the structure given a path to a sav file
func (sav *saveGameType) loadSave(path string) error {
	var err error
	sav.data, _, err = civ3decompress.ReadFile(path)
	if err != nil {
		return err
	}
	sav.path = path
	sav.populateSections()
	// If this is a save game, populate BIQ and save vars
	// This is still hackish
	if string(sav.data[0:4]) == "CIV3" {
		gameOff, err := sav.sectionOffset("GAME", 2)
		if err != nil {
			return nil
		}
		currentGame.data = sav.data[gameOff-4:]
		currentGame.populateSections()
		bicOff, err := sav.sectionOffset("VER#", 1)
		if err != nil {
			return nil
		}
		if sav.readInt32(bicOff+8, Unsigned) == 0xcdcdcdcd {
			currentBic = defaultBic
		} else {
			currentBic.data = sav.data[bicOff-8 : gameOff]
			currentBic.populateSections()
		}
		currentGame.data = saveGame.data[gameOff-4:]
	}
	return nil
}

// Find sections demarc'ed by 4-character ASCII headers and place into sections[]
func (sav *saveGameType) populateSections() {
	var i, count, offset int
	sav.sections = make([]sectionType, 0)
	// find sections demarc'ed by 4-character ASCII headers
	for i < len(sav.data) {
		if sav.data[i] < 0x20 || sav.data[i] > 0x5a {
			count = 0
		} else {
			if count == 0 {
				offset = i
			}
			count++
		}
		i++
		if count > 3 {
			count = 0
			s := new(sectionType)
			s.offset = offset
			s.name = string(sav.data[offset:i])
			sav.sections = append(sav.sections, *s)
		}
	}
}

// returns just the filename part of the path assuming / or \ separators
func (sav *saveGameType) fileName() string {
	var o int
	for i := 0; i < len(sav.path); i++ {
		if sav.path[i] == 0x2f || sav.path[i] == 0x5c {
			o = i
		}
	}
	return sav.path[o+1:]
}

func (sav *saveGameType) sectionOffset(sectionName string, nth int) (int, error) {
	var i, n int
	for i < len(sav.sections) {
		if sav.sections[i].name == sectionName {
			n++
			if n >= nth {
				return sav.sections[i].offset + len(sectionName), nil
			}
		}
		i++
	}
	return -1, errors.New("Could not find " + strconv.Itoa(nth) + " section named " + sectionName)
}

func (sav *saveGameType) readInt32(offset int, signed bool) int {
	n := int(sav.data[offset]) +
		int(sav.data[offset+1])*0x100 +
		int(sav.data[offset+2])*0x10000 +
		int(sav.data[offset+3])*0x1000000
	if signed && n&0x80000000 != 0 {
		n = -(n ^ 0xffffffff + 1)
	}
	return n
}

func (sav *saveGameType) readInt16(offset int, signed bool) int {
	n := int(sav.data[offset]) +
		int(sav.data[offset+1])*0x100
	if signed && n&0x8000 != 0 {
		n = -(n ^ 0xffff + 1)
	}
	return n
}

func (sav *saveGameType) readInt8(offset int, signed bool) int {
	n := int(sav.data[offset])
	if signed && n&0x80 != 0 {
		n = -(n ^ 0xff + 1)
	}
	return n
}
