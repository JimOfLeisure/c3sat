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

// Flavor is the leaf element of FLAV
type Flavor struct {
	A                      int32
	FlavorName             [0x100]byte
	B, C, D, E, F, G, H, I int32
}

// Game is the first section after the BIC.
type Game struct {
	// First two fields count for "class base"
	Name                       [4]byte
	_                          int32
	_                          [3]int32
	RenderFlags                int32
	DifficultyLevel            int32
	_                          int32
	UnitsCount                 int32
	CitiesCount                int32
	_                          int32
	_                          int32
	GlobalWarmingLevel         int32
	_                          int32
	_                          int32
	_                          int32
	CurrentTurn                int32
	_                          int32
	Random                     int32
	_                          int32
	CivFlags2                  int32
	CivFlags1                  int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          int32
	_                          [48]int32
	Value1                     int32
	_                          [72]int32
	GameLimitPoints            int32
	GameLimitTurns             int32
	_                          [50]int32
	_                          int32
	_                          int32
	GameLimitDestroyedCities   int32
	GameLimitCityCulture       int32
	GameLimitCivCulture        int32
	GameLimitPopulation        int32
	GameLimitTerritory         int32
	GameLimitWonders           int32
	GameLimitDestroyedWonders  int32
	GameLimitAdvances          int32
	GameLimitCapturedCities    int32
	GameLimitVictoryPointPrice int32
	GameLimitPrincessRansom    int32
	DefaultDate1               int32
	_                          [27]int32
	PLGI                       [10]int32
	Date2                      Date
	Date3                      Date
	GameAggression             int32
	_                          int32
	CityStatIntArray           int32
	ResearchedAdvances         int32
	Wonders                    int32
	WonderFlags                int32
	ImprovementTypesData1      int32
	ImprovementTypesData2      int32
	_                          int32
	_                          int32
}

// Date DATE section
type Date struct {
	Name         [4]byte
	Length       int32
	Text         [16]byte
	_            [12]int32
	BaseTimeUnit int32
	Month        int32
	Week         int32
	Year         int32
	_            int32
}

// MyTest ...
type MyTest struct {
	Name [4]byte
}
