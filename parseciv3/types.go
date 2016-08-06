package parseciv3

// Civ3Data contains the game data
type Civ3Data struct {
	FileName   string
	Compressed bool
	RawFile    []byte
}

type civ3 struct {
	Name [4]byte
	// 28 bytes. Guessing on alignment
	A, B, C, D, E, F uint32
	G                uint16
}
