package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the PRTO sectionsof the BIC
func prtoModule(L *lua.LState) {
	prto := L.NewTable()
	L.SetGlobal("prto", prto)
	prtoOff, _ := currentBic.sectionOffset("PRTO", 1)
	listSection(prtoOff, func(off int, length int) {
		lt := L.NewTable()
		prto.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off+4:off+4+32])))
		L.RawSet(lt, lua.LString("attack"), lua.LNumber(currentBic.readInt32(off+92, Signed)))
		L.RawSet(lt, lua.LString("defense"), lua.LNumber(currentBic.readInt32(off+84, Signed)))
		L.RawSet(lt, lua.LString("move"), lua.LNumber(currentBic.readInt32(off+108, Signed)))
		L.RawSet(lt, lua.LString("cost"), lua.LNumber(currentBic.readInt32(off+80, Signed)))
		L.RawSet(lt, lua.LString("transport"), lua.LNumber(currentBic.readInt32(off+76, Signed)))
	})
}
