package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the RACE sectionsof the BIC
func raceModule(L *lua.LState) {
	// var name string
	race := L.NewTable()
	L.SetGlobal("race", race)
	raceOff, _ := currentBic.sectionOffset("RACE", 1)
	listSection(raceOff, func(off int, length int) {
		lt := L.NewTable()
		race.Append(lt)
		numCityNames := currentBic.readInt32(off, Signed)
		off += 4
		cityNames := L.NewTable()
		L.RawSet(lt, lua.LString("city_names"), cityNames)
		for i := 0; i < numCityNames; i++ {
			cityNames.Append(lua.LString(civString(currentBic.data[off : off+24])))
			off += 24
		}
		numGreatLeaders := currentBic.readInt32(off, Signed)
		off += 4
		greatLeaders := L.NewTable()
		L.RawSet(lt, lua.LString("great_leader_names"), greatLeaders)
		for i := 0; i < numGreatLeaders; i++ {
			greatLeaders.Append(lua.LString(civString(currentBic.data[off : off+32])))
			off += 32
		}
		L.RawSet(lt, lua.LString("leader_name"), lua.LString(civString(currentBic.data[off:off+32])))
		off += 32
		L.RawSet(lt, lua.LString("leader_title"), lua.LString(civString(currentBic.data[off:off+24])))
		off += 24
		// skip over civilopedia entry
		off += 32
		L.RawSet(lt, lua.LString("adjective"), lua.LString(civString(currentBic.data[off:off+40])))
		off += 40
		L.RawSet(lt, lua.LString("name"), lua.LString(civString(currentBic.data[off:off+40])))
		off += 40
		L.RawSet(lt, lua.LString("object_noun"), lua.LString(civString(currentBic.data[off:off+40])))
		off += 40
		// non-string info here, and then list of scientific leaders at the end
	})
}
