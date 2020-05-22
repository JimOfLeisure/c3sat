package luaciv3

import (
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
		L.RawSet(lt, lua.LString("content_citizens"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		L.RawSet(lt, lua.LString("max_anarchy_turns"), lua.LNumber(currentBic.readInt32(off+68, Signed)))
		L.RawSet(lt, lua.LString("defense_land_units"), lua.LNumber(currentBic.readInt32(off+72, Signed)))
		L.RawSet(lt, lua.LString("offense_land_units"), lua.LNumber(currentBic.readInt32(off+76, Signed)))
		L.RawSet(lt, lua.LString("start_units_1"), lua.LNumber(currentBic.readInt32(off+80, Signed)))
		L.RawSet(lt, lua.LString("start_units_2"), lua.LNumber(currentBic.readInt32(off+84, Signed)))
		L.RawSet(lt, lua.LString("add_free_support"), lua.LNumber(currentBic.readInt32(off+88, Signed)))
		L.RawSet(lt, lua.LString("bonus_each_city"), lua.LNumber(currentBic.readInt32(off+92, Signed)))
		L.RawSet(lt, lua.LString("barb_attack_bonus"), lua.LNumber(currentBic.readInt32(off+96, Signed)))
		L.RawSet(lt, lua.LString("cost_factor"), lua.LNumber(currentBic.readInt32(off+100, Signed)))
		L.RawSet(lt, lua.LString("pct_optimal_cities"), lua.LNumber(currentBic.readInt32(off+104, Signed)))
		L.RawSet(lt, lua.LString("ai_trade_rate"), lua.LNumber(currentBic.readInt32(off+108, Signed)))
		L.RawSet(lt, lua.LString("corruption"), lua.LNumber(currentBic.readInt32(off+112, Signed)))
		L.RawSet(lt, lua.LString("quelled_citizens"), lua.LNumber(currentBic.readInt32(off+116, Signed)))
		// L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
	})
}
