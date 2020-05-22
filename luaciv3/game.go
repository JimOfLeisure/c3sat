package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

func gameModule(L *lua.LState) {
	game := L.NewTable()
	L.SetGlobal("game", game)
	gameOff, _ := currentGame.sectionOffset("GAME", 1)
	L.RawSet(game, lua.LString("city_count"), lua.LNumber(currentGame.readInt32(gameOff+32, Signed)))
	L.RawSet(game, lua.LString("unit_count"), lua.LNumber(currentGame.readInt32(gameOff+28, Signed)))
}
