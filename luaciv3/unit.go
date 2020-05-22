package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the LEAD sections which are the 32 players in a game
// Player 0 is the barbarians, player 1 is the first human player
func unitModule(L *lua.LState) {
	// const unitLen = 536
	unit := L.NewTable()
	L.SetGlobal("unit", unit)
	for _, v := range currentGame.sections {
		if v.name == "UNIT" {
			thisUnit := L.NewTable()
			unit.Append(thisUnit)
			// this seems to include the 4-char header
			offset := v.offset + 4
			L.RawSet(thisUnit, lua.LString("id"), lua.LNumber(currentGame.readInt32(offset+4, Signed)))
			L.RawSet(thisUnit, lua.LString("x"), lua.LNumber(currentGame.readInt32(offset+8, Signed)))
			L.RawSet(thisUnit, lua.LString("y"), lua.LNumber(currentGame.readInt32(offset+12, Signed)))
			L.RawSet(thisUnit, lua.LString("prev_x"), lua.LNumber(currentGame.readInt32(offset+16, Signed)))
			L.RawSet(thisUnit, lua.LString("prev_y"), lua.LNumber(currentGame.readInt32(offset+20, Signed)))
			L.RawSet(thisUnit, lua.LString("civ_id"), lua.LNumber(currentGame.readInt32(offset+24, Signed)))
			L.RawSet(thisUnit, lua.LString("race_id"), lua.LNumber(currentGame.readInt32(offset+28, Signed)))
			L.RawSet(thisUnit, lua.LString("whats_this"), lua.LNumber(currentGame.readInt32(offset+32, Signed)))
			// I thought this might be prto_id, but it doesn't seem to be. Antal1987 calls it UnitTypeID
			L.RawSet(thisUnit, lua.LString("unit_type_id"), lua.LNumber(currentGame.readInt32(offset+36, Signed)))
			// for i := 0; i < 32; i++ {
			// 	thisUnit.Append(lua.LNumber(currentGame.readInt32(offset+i*4, Signed)))
			// }
		}
	}
}
