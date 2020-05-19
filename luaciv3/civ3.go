package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

func civ3Module(L *lua.LState) {
	civ3 := L.NewTable()
	L.SetGlobal("civ3", civ3)
}
