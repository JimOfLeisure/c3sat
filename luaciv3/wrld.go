package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

func wrldModule(L *lua.LState) {
	wrld := L.NewTable()
	L.SetGlobal("wrld", wrld)
	offset, _ := currentGame.sectionOffset("WRLD", 1)
	L.RawSet(wrld, lua.LString("wsiz_id"), lua.LNumber(currentGame.readInt32(offset+234, Signed)))
}
