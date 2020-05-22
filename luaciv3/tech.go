package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the TECH sectionsof the BIC
func techModule(L *lua.LState) {
	tech := L.NewTable()
	L.SetGlobal("tech", tech)
	techOff, _ := currentBic.sectionOffset("TECH", 1)
	listSection(techOff, func(off int, length int) {
		lt := L.NewTable()
		tech.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off:off+32])))
	})
}
