package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the GOVT sectionsof the BIC
func govtModule(L *lua.LState) {
	govt := L.NewTable()
	L.SetGlobal("govt", govt)
	govtOff, _ := currentBic.sectionOffset("GOVT", 1)
	listSection(govtOff, func(off int) {
		lt := L.NewTable()
		govt.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off+24:off+24+64])))
	})
}
