package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the PRTO sectionsof the BIC
func prtoModule(L *lua.LState) {
	var prtoLen int
	prto := L.NewTable()
	L.SetGlobal("prto", prto)
	prtoOff, _ := currentBic.sectionOffset("PRTO", 1)
	numPrto := currentBic.readInt32(prtoOff, Signed)
	off := prtoOff + 4
	for i := 0; i < numPrto; i++ {
		lt := L.NewTable()
		prto.Append(lt)
		prtoLen = currentBic.readInt32(off, Signed)
		// skip over the length
		off += 4
		name, err := CivString(currentBic.data[off+4 : off+4+32])
		if err != nil {
			// TODO: handle errors
			panic(err)
		}
		L.RawSet(lt, lua.LString("name"), lua.LString(name))
		L.RawSet(lt, lua.LString("attack"), lua.LNumber(currentBic.readInt32(off+92, Signed)))
		L.RawSet(lt, lua.LString("defense"), lua.LNumber(currentBic.readInt32(off+84, Signed)))
		L.RawSet(lt, lua.LString("move"), lua.LNumber(currentBic.readInt32(off+108, Signed)))
		L.RawSet(lt, lua.LString("cost"), lua.LNumber(currentBic.readInt32(off+80, Signed)))
		L.RawSet(lt, lua.LString("transport"), lua.LNumber(currentBic.readInt32(off+76, Signed)))
		off += prtoLen
	}
}
