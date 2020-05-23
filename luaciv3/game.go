package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

func gameModule(L *lua.LState) {
	game := L.NewTable()
	L.SetGlobal("game", game)
	gameOff, _ := currentGame.sectionOffset("GAME", 1)
	L.RawSet(game, lua.LString("diff_id"), lua.LNumber(currentGame.readInt32(gameOff+20, Signed)))
	L.RawSet(game, lua.LString("unit_count"), lua.LNumber(currentGame.readInt32(gameOff+28, Signed)))
	L.RawSet(game, lua.LString("city_count"), lua.LNumber(currentGame.readInt32(gameOff+32, Signed)))
	// The per-civ tech list is actually a per-tech 32-bit bitmask, and the number of continents impacts its offset
	techs := L.NewTable()
	L.RawSet(game, lua.LString("tech_civ_bitmask"), techs)
	techOff, _ := currentBic.sectionOffset("TECH", 1)
	numTechs := currentBic.readInt32(techOff, Signed)
	wrldOff, _ := currentGame.sectionOffset("WRLD", 1)
	continentCount := currentGame.readInt16(wrldOff+4, Signed)
	techCivMaskOffset := 856 + 4*continentCount
	for i := 0; i < numTechs; i++ {
		techs.Append(lua.LNumber(currentGame.readInt32(techCivMaskOffset+4*i, Unsigned)))
	}
}
