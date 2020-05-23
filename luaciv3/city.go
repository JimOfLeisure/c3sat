package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

func cityModule(L *lua.LState) {
	var count int
	city := L.NewTable()
	L.SetGlobal("city", city)
	for _, v := range currentGame.sections {
		if v.name == "CITY" {
			count++
		}
		L.RawSet(city, lua.LString("count"), lua.LNumber(count))
	}
}
