package luaciv3

import (
	"fmt"

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
	fmt.Println(numPrto)
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
		off += prtoLen
	}
}
