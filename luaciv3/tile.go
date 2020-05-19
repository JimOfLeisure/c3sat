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
	Lua table arrays usually start with index 1, but this starts at 0
	Only use ipairs when parsing the tiles because I'll add string keys for
	useful info
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
}
