package parseciv3

import "fmt"

// Debug ...
func (c Civ3Data) Debug() string {
	var out string
	section := []string{"LEAD", "GAME", "GAME2"}
	for v := range section {
		// out += fmt.Sprintln(section[v])
		if civ3Data, ok := c.Data[section[v]].(List); ok {
			out += fmt.Sprintf("%s Name: %s Count: %x\n", section[v], civ3Data.Name, civ3Data.Count)
		}
		if civ3Data, ok := c.Data[section[v]].(Base); ok {
			out += fmt.Sprintf("%s Name: %s Length: %x\n", section[v], civ3Data.Name, civ3Data.Length)
		}
	}
	// out += fmt.Sprintf("\nLEAD Name: %s Length: %x\n", c.Data["LEAD"].Name, c.Data["LEAD"].Length)
	// out += fmt.Sprintf("\GAME Name: %s Length: %x\n", c.Data["GAME"].Name, c.Data["GAME"].Length)
	// out += fmt.Sprintf("*** GAME ***\n%#v\n", c.Data["GAME2"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["GameNext"])
	// out += fmt.Sprintf("\n%s\n", c.Data["WTF"])
	out += fmt.Sprintf("\n%#v\n", c.Data["WRLD"])

	out += fmt.Sprint("*** Debug output. Next bytes ***\n")
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
