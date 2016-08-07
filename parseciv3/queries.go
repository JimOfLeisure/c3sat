package parseciv3

import "fmt"

// Debug ...
func (c Civ3Data) Debug() string {
	var out string
	out += fmt.Sprintf("\n%v\n", c.Data["LEAD"])
	// out += fmt.Sprintf("\n*** GAME ***\n%s\n", hex.Dump(c.Data["GAME"].data))
	out += fmt.Sprint("\n*** Debug output. Next bytes ***\n\n")
	out += fmt.Sprintln(c.Next)
	return out
}

// Info ...
func (c Civ3Data) Info() string {
	var out string
	out += fmt.Sprintf("File: %s\t", c.FileName)
	out += fmt.Sprintf("Compressed: %v\n", c.Compressed)
	return out
}
