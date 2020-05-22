package luaciv3

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Provides data from the PRTO sectionsof the BIC
func prtoModule(L *lua.LState) {
	// hack
	const prtoLen = 259
	prto := L.NewTable()
	L.SetGlobal("prto", prto)
	prtoOff, _ := currentBic.sectionOffset("PRTO", 1)
	numPrto := currentBic.readInt32(prtoOff, Signed)
	fmt.Println(numPrto)
	for i := 0; i < numPrto; i++ {
		lt := L.NewTable()
		prto.Append(lt)
		off := prtoOff + 12 + (i * prtoLen)
		// FIXME
		// this works until after Submarine ... ??? [36] or in lua [37]
		// Maybe the first offset is an int length not counting its own data; that adds up
		name, err := CivString(currentBic.data[off : off+32])
		if err != nil {
			// TODO: handle errors
			panic(err)
		}
		L.RawSet(lt, lua.LString("name"), lua.LString(name))
	}
}
