package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

/*
	CIV3 is the 'section string' at the beginning of every Civ III save
	I don't really understand what's in there, but below are my best guesses
	so far

	This creates the global 'civ3' table in Lua with values from the last sav.load()
*/

func civ3Module(L *lua.LState) {
	civ3 := L.NewTable()
	L.SetGlobal("civ3", civ3)
	L.RawSet(civ3, lua.LString("always26"), lua.LNumber(saveGame.readInt8(5, Signed)))
	L.RawSet(civ3, lua.LString("maybe_version_minor"), lua.LNumber(saveGame.readInt32(6, Signed)))
	L.RawSet(civ3, lua.LString("maybe_version_major"), lua.LNumber(saveGame.readInt32(10, Signed)))
	L.RawSet(civ3, lua.LString("gobbledegook_1"), lua.LNumber(saveGame.readInt32(14, Unsigned)))
	L.RawSet(civ3, lua.LString("gobbledegook_2"), lua.LNumber(saveGame.readInt32(18, Unsigned)))
	L.RawSet(civ3, lua.LString("gobbledegook_3"), lua.LNumber(saveGame.readInt32(22, Unsigned)))
	L.RawSet(civ3, lua.LString("gobbledegook_4"), lua.LNumber(saveGame.readInt32(26, Unsigned)))
}
