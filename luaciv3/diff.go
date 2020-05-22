package luaciv3

import (
	"encoding/hex"

	lua "github.com/yuin/gopher-lua"
)

// Provides data from the DIFF sectionsof the BIC
func diffModule(L *lua.LState) {
	diff := L.NewTable()
	L.SetGlobal("diff", diff)
	offset, _ := currentBic.sectionOffset("DIFF", 1)
	listSection(offset, func(off int, length int) {
		lt := L.NewTable()
		diff.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off:off+64])))
		L.RawSet(lt, lua.LString("0"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		L.RawSet(lt, lua.LString("1"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		L.RawSet(lt, lua.LString("2"), lua.LNumber(currentBic.readInt32(off+72, Signed)))
		L.RawSet(lt, lua.LString("3"), lua.LNumber(currentBic.readInt32(off+76, Signed)))
		L.RawSet(lt, lua.LString("4"), lua.LNumber(currentBic.readInt32(off+80, Signed)))
		L.RawSet(lt, lua.LString("5"), lua.LNumber(currentBic.readInt32(off+84, Signed)))
		L.RawSet(lt, lua.LString("6"), lua.LNumber(currentBic.readInt32(off+88, Signed)))
		L.RawSet(lt, lua.LString("7"), lua.LNumber(currentBic.readInt32(off+92, Signed)))
		L.RawSet(lt, lua.LString("8"), lua.LNumber(currentBic.readInt32(off+96, Signed)))
		L.RawSet(lt, lua.LString("9"), lua.LNumber(currentBic.readInt32(off+100, Signed)))
		L.RawSet(lt, lua.LString("10"), lua.LNumber(currentBic.readInt32(off+104, Signed)))
		L.RawSet(lt, lua.LString("11"), lua.LNumber(currentBic.readInt32(off+108, Signed)))
		L.RawSet(lt, lua.LString("12"), lua.LNumber(currentBic.readInt32(off+112, Signed)))
		L.RawSet(lt, lua.LString("13"), lua.LNumber(currentBic.readInt32(off+116, Signed)))
		L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
	})
}
