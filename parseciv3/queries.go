package parseciv3

import (
	"encoding/json"
	"fmt"
)

// Debug ...
func (c Civ3Data) Debug() string {
	var out string
	// section := []string{"LEAD", "GAME", "GAME2", "GOOD"}
	// for v := range section {
	// 	// out += fmt.Sprintln(section[v])
	// 	if civ3Data, ok := c.Data[section[v]].(List); ok {
	// 		out += fmt.Sprintf("%s Name: %s Count: %x\n", section[v], civ3Data.Name, civ3Data.Count)
	// 	}
	// 	if civ3Data, ok := c.Data[section[v]].(Base); ok {
	// 		out += fmt.Sprintf("%s Name: %s Length: %x\n", section[v], civ3Data.Name, civ3Data.Length)
	// 	}
	// }
	// out += fmt.Sprintf("\nLEAD Name: %s Length: %x\n", c.Data["LEAD"].Name, c.Data["LEAD"].Length)
	// out += fmt.Sprintf("\GAME Name: %s Length: %x\n", c.Data["GAME"].Name, c.Data["GAME"].Length)
	// out += fmt.Sprintf("*** GAME ***\n%#v\n", c.Data["GAME2"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["GameNext"])
	// out += fmt.Sprintf("\n%s\n", c.Data["WTF"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["WRLD"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["CNSL"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["TILE"])
	// out += fmt.Sprintf("\n%#v\n", c.Data["ResourceCounts"])

	// out += fmt.Sprint("*** Debug output. Next bytes ***\n")
	// out += fmt.Sprintln(c.Next)
	// out += fmt.Sprintln(json.Marshal(string(c.Data["TILE"])))
	return out
}

// Info ...
func (c Civ3Data) Info() string {
	var out string
	out += fmt.Sprintf("File: %s\t", c.FileName)
	out += fmt.Sprintf("Compressed: %v\n", c.Compressed)
	return out
}

// Map ...
type Map struct {
	Width                  int32  `json:"width"`
	Height                 int32  `json:"height"`
	Tile                   []Tile `json:"tile"`
	CivStartLocationTileID [32]int32
}

// JSONMap returns byte array to write to JSON file, then read with html/d3.html map reader in this repo
func (c Civ3Data) JSONMap() []byte {
	var civmap Map
	if tiles, ok := c.Data["TILE"].([]Tile); ok {
		civmap.Tile = tiles
	}
	if wrld, ok := c.Data["WRLD"].(Wrld); ok {
		civmap.Height = wrld.MapHeight
		civmap.Width = wrld.MapWidth
		civmap.CivStartLocationTileID = wrld.CivStartLocationTileID
	}
	out, _ := json.Marshal(civmap)
	return out
}
