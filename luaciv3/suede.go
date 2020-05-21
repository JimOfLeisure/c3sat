package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Was unsure what to implement next in the project, got a little inspiration from
//  https://forums.civfanatics.com/threads/when-do-tier-2-barbarian-units-spawn.657845/
//  So I'll try to count cities, barb units and types, and barb camps
//  Eventually this will be refactored to other modules, but for now we'll call it suede after OP
func suedeModule(L *lua.LState) {
	suede := L.NewTable()
	L.SetGlobal("suede", suede)
	gameOff, _ := saveGame.sectionOffset("GAME", 2)
	L.RawSet(suede, lua.LString("city_count"), lua.LNumber(saveGame.readInt32(gameOff+32, Signed)))
}
