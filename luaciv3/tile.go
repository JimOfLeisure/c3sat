package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

/*
	TILE - There is one TILE value for each tile on the map. Actually there are
	four 'TILE' 'sections' in the save file per tile, but we'll expose the tiles
	individually.

	Note that one tile left/right is x-2/x+2, up/down is y-2/y+2, and diagonal
	is +/-1 in each axis. Even y rows start with x=0, and odd y rows start with x=1.

	This creates the global 'tile' table in Lua with values from the last sav.load()
	Only use ipairs when parsing the tiles because I'll add string keys for
	useful info

	Note that save file array indexes start with 0 and Lua array indexes start with 1
*/

func tileModule(L *lua.LState) {
	tile := L.NewTable()
	L.SetGlobal("tile", tile)
	// Get 2nd WRLD offset
	section, err := saveGame.sectionOffset("WRLD", 2)
	if err != nil {
		panic(err)
	}
	// Reading 6 int32s at offset 8; first is height, last is Width
	intList := make([]int, 6)
	for i := 0; i < 6; i++ {
		intList[i] = saveGame.readInt32(section+8+4*i, Signed)
	}

	L.RawSet(tile, lua.LString("height"), lua.LNumber(intList[0]))
	L.RawSet(tile, lua.LString("width"), lua.LNumber(intList[5]))

	const tileBytes = 212
	tilesOffset, err := saveGame.sectionOffset("TILE", 1)
	var mapRowLength = intList[5] / 2
	var mapTileCount = mapRowLength * intList[0]
	// mapTileOffsets := make([]int, mapTileCount)
	for i := 0; i < mapTileCount; i++ {
		thisTile := L.NewTable()
		tile.Append(thisTile)
		// L.RawSet(tile, lua.LNumber(i), thisTile)
		tileOffset := tilesOffset - 4 + (i/mapRowLength)*mapRowLength*tileBytes + (i%mapRowLength)*tileBytes
		// mapTileOffsets[i] = tileOffset
		terrain := saveGame.readInt8(tileOffset+57, Unsigned)
		L.RawSet(thisTile, lua.LString("terrain"), lua.LNumber(terrain))
		L.RawSet(thisTile, lua.LString("base_terrain"), lua.LNumber(terrain&0x0f))
		L.RawSet(thisTile, lua.LString("overlay_terrain"), lua.LNumber(terrain>>4))
	}
}
