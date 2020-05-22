package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the WSIZ sectionsof the BIC
func wsizModule(L *lua.LState) {
	wsiz := L.NewTable()
	L.SetGlobal("wsiz", wsiz)
	offset, _ := currentBic.sectionOffset("WSIZ", 1)
	listSection(offset, func(off int) {
		lt := L.NewTable()
		wsiz.Append(lt)
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off+32:off+32+32])))
		// No idea if this is byte, short, or int, or if it's signed; Antal1987 suggests it's int
		L.RawSet(lt, lua.LString("ocn"), lua.LNumber(currentBic.readInt32(off, Signed)))
		L.RawSet(lt, lua.LString("width"), lua.LNumber(currentBic.readInt32(off+64, Signed)))
		L.RawSet(lt, lua.LString("dist_between_civs"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		// TODO: it's possible I have height and width backwards; all the default sizes are square; check on this
		L.RawSet(lt, lua.LString("num_civs"), lua.LNumber(currentBic.readInt32(off+72, Signed)))
		L.RawSet(lt, lua.LString("height"), lua.LNumber(currentBic.readInt32(off+76, Signed)))
		L.RawSet(lt, lua.LString("whats_this2"), lua.LNumber(currentBic.readInt32(off+80, Signed)))
		L.RawSet(lt, lua.LString("whats_this3"), lua.LNumber(currentBic.readInt32(off+84, Signed)))
		// L.RawSet(lt, lua.LString("dump"), lua.LString(hex.Dump(currentBic.data[off:off+256])))
	})
}
