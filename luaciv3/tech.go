package luaciv3

import (
	"encoding/hex"

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
		L.RawSet(lt, lua.LString("era_id"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		prereq := L.NewTable()
		L.RawSet(lt, lua.LString("prereq_tech_ids"), prereq)
		for i := 0; i < 4; i++ {
			pre := currentBic.readInt32(off+84+i*4, Signed)
			if pre >= 0 {
				prereq.Append(lua.LNumber(pre))
			}
		}
		L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
	})
}
