package parseciv3

import (
	"encoding/hex"
	"fmt"
)

// Debug ...
func (c Civ3Data) Debug() string {
	var out string
	out += fmt.Sprintln("Debug output:")
	out += fmt.Sprintln(hex.Dump(c.RawFile[:0x100]))
	return out
}

// Info ...
func (c Civ3Data) Info() string {
	var out string
	out += fmt.Sprintf("File: %s\t", c.FileName)
	out += fmt.Sprintf("Compressed: %v\n", c.Compressed)
	return out
}
