package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the LEAD sections which are the 32 players in a game
// Player 0 is the barbarians, player 1 is the first human player
func unitModule(L *lua.LState) {
	const unitLen = 536
	var unitOff int
	var offset int
	unit := L.NewTable()
	L.SetGlobal("unit", unit)
	// Since "UNIT" can and does often appear in dirty data in unitialized parts of the save,
	// will search for UNIT after that
	lastEspn, _ := currentGame.sectionOffset("ESPN", 64)
	for i := 0; i < 50; i++ {
		unitOff, _ = currentGame.sectionOffset("UNIT", i+1)
		if unitOff > lastEspn {
			break
		}
	}
	gameOff, _ := currentGame.sectionOffset("GAME", 1)
	numUnits := currentGame.readInt32(gameOff+28, Signed)
	offset = unitOff
	// for _, v := range currentGame.sections {
	for i := 0; i < numUnits; i++ {
		// if v.name == "UNIT" {
		thisUnit := L.NewTable()
		unit.Append(thisUnit)
		// this seems to include the 4-char header
		// offset := v.offset + 4
		L.RawSet(thisUnit, lua.LString("id"), lua.LNumber(currentGame.readInt32(offset+4, Signed)))
		L.RawSet(thisUnit, lua.LString("x"), lua.LNumber(currentGame.readInt32(offset+8, Signed)))
		L.RawSet(thisUnit, lua.LString("y"), lua.LNumber(currentGame.readInt32(offset+12, Signed)))
		L.RawSet(thisUnit, lua.LString("prev_x"), lua.LNumber(currentGame.readInt32(offset+16, Signed)))
		L.RawSet(thisUnit, lua.LString("prev_y"), lua.LNumber(currentGame.readInt32(offset+20, Signed)))
		L.RawSet(thisUnit, lua.LString("civ_id"), lua.LNumber(currentGame.readInt32(offset+24, Signed)))
		L.RawSet(thisUnit, lua.LString("race_id"), lua.LNumber(currentGame.readInt32(offset+28, Signed)))
		// L.RawSet(thisUnit, lua.LString("whats_this"), lua.LNumber(currentGame.readInt32(offset+32, Signed)))
		L.RawSet(thisUnit, lua.LString("prto_id"), lua.LNumber(currentGame.readInt32(offset+36, Signed)))
		offset += unitLen
	}
}
