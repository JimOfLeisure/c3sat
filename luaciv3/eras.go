package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the ERAS sectionsof the BIC
func erasModule(L *lua.LState) {
	eras := L.NewTable()
	L.SetGlobal("eras", eras)
	erasOff, _ := currentBic.sectionOffset("ERAS", 1)
	listSection(erasOff, func(off int) {
		lt := L.NewTable()
		eras.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off:off+64])))
	})
}
