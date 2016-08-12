package parseciv3

import (
	"encoding/binary"
	"io"
)

/*
The save file seems to be a simple binary dump of C++ data structures
packed with no byte padding. Generally speaking, most of the data is in
classes inherited from two basic classes. Both start with a 4-byte string
which appears to be a class name closely related to its function. One class
then has a 32-bit integer expressing the length in bytes of the data structure
following. The other has a 32-bit integer as a count of records. Each record
begins with a 32-bit length in bytes followed by the data. Before I knew this
I called each labeled length a "section", so I'll sometimes use that term
even now.

Some non-conformers appear to be the inital CIV3 section, but it's at least a
consistent length. The FLAV section is a list of lists. The second GAME section
in the SAV (which is the first GAME section of the non-BIC info) has an
apparently meaningless integer after the header followed by some predictable
data and then some as-yet unpredictable data which may be integer arrays, but
I haven't yet found the count. The length in bytes from GAME to DATE seems to
always be odd, so there must be a lone byte or a byte array in there somewhere.
I found a couple of stray apparent int32s after one of the DATE sections.

Later after the map tile data I have yet to figure out, too.

My strategy in the 2013-2015 Python version of this parser and my strategy so
far in Go is to parse the header, length/count and the data and then interpret
it. But several of the sections repeat with different lengths and data, especially
TILE but also WRLD and some others. I am presuming this is due to successive
versions of the game inheriting classes from the earlier game and adding to them,
and it shows up in the SAV file as the inheritance chain with data from each
generation. Mechanically parsing lenghts and counts works, but there really is
no advantage in meaning.

So I'm going to instead start making Go structs that will capture the entire
inheritance chain in one read which should make more sense programatically.

As I type this, I am mechanically parsing to WRLD but extracting little meaning
so far. During transition I'll be reading with two different methods.
*/

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
}

// GameNext is what Antal1987's dumps suggest is next, but I don't think so
type GameNext struct {
	_                     [27]int32
	PLGI                  [10]int32
	Date2                 Date
	Date3                 Date
	GameAggression        int32
	_                     int32
	CityStatIntArray      int32
	ResearchedAdvances    int32
	Wonders               int32
	WonderFlags           int32
	ImprovementTypesData1 int32
	ImprovementTypesData2 int32
	UnitTypesData1        int32
	UnitTypesData2        int32
	_                     int32
	_                     int32
	DefaultGameSettings   DefaultGameSettings
}

// Date DATE section ... I don't think this is nearly right
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

// DefaultGameSettings from Antal1987's dump
type DefaultGameSettings struct {
	TurnsLimit           int32
	PointsLimit          int32
	DestroyedCitiesLimit int32
	CityCultureLimit     int32
	CivCultureLimit      int32
	PopulationLimit      int32
	TerritoryLimit       int32
	WondersLimit         int32
	DestroyedUnitsLimit  int32
	AdvancesLimit        int32
	CapturedCitiesLimit  int32
	VictoryPointPrice    int32
	PrincessPrice        int32
	PrincessRansom       int32
}

// MyTest ...
type MyTest struct {
	Name [4]byte
}
