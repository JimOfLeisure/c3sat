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
	gameOff, _ := currentGame.sectionOffset("GAME", 1)
	// This appears to be plausible for the global city count, and unit count
	L.RawSet(suede, lua.LString("city_count"), lua.LNumber(currentGame.readInt32(gameOff+32, Signed)))
	L.RawSet(suede, lua.LString("unit_count"), lua.LNumber(currentGame.readInt32(gameOff+28, Signed)))
	// L.RawSet(suede, lua.LString("sections"), lua.LNumber(currentGame.readInt32(gameOff+28, Signed)))
	var count int
	foo := L.NewTable()
	L.RawSet(suede, lua.LString("sizes"), foo)
	var lastOff int
	for _, v := range currentGame.sections {
		if v.name == "UNIT" {
			if lastOff != 0 {
				foo.Append(lua.LNumber(v.offset - lastOff))
			}
			lastOff = v.offset
			count++
		}
	}
	L.RawSet(suede, lua.LString("unit_sections"), lua.LNumber(count))
}
