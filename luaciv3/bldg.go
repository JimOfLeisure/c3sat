package luaciv3

import (
	"encoding/hex"

	lua "github.com/yuin/gopher-lua"
)

// Provides data from the BLDG sectionsof the BIC
func bldgModule(L *lua.LState) {
	bldg := L.NewTable()
	L.SetGlobal("bldg", bldg)
	offset, _ := currentBic.sectionOffset("BLDG", 1)
	listSection(offset, func(off int, length int) {
		lt := L.NewTable()
		bldg.Append(lt)
		L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
	})
}
