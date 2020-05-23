package luaciv3

import (
	"encoding/hex"

	lua "github.com/yuin/gopher-lua"
)

func cityModule(L *lua.LState) {
	// var count int
	var offset int
	lt := L.NewTable()
	L.SetGlobal("city", lt)
	// Since "CITY" can sometimes appear in dirty data, let's find the last IDLS
	//   and find the first "CITY" after that
	var lastIdls int
	for _, v := range currentGame.sections {
		if v.name == "IDLS" {
			lastIdls = v.offset
		}
	}
	for _, v := range currentGame.sections {
		if v.name == "CITY" && v.offset > lastIdls {
			offset = v.offset
			break
		}
	}
	// Get city count
	gameOff, _ := currentGame.sectionOffset("GAME", 1)
	numCities := currentGame.readInt32(gameOff+28, Signed)

	// TODO: if offset is still 0 we have an error... or we have 0 cities
	for i := 0; i < numCities; i++ {
		thisCity := L.NewTable()
		lt.Append(thisCity)
		// offsets here are from the start of "CITY"
		// offset is right, but not sure about length; 20 is what Antal1987's dumps say
		L.RawSet(thisCity, lua.LString("name"), lua.LString(civString(currentGame.data[offset+0x188:offset+0x188+20])))
		L.RawSet(thisCity, lua.LString("id"), lua.LNumber(currentGame.readInt32(offset+8, Signed)))
		L.RawSet(thisCity, lua.LString("x"), lua.LNumber(currentGame.readInt16(offset+12, Signed)))
		L.RawSet(thisCity, lua.LString("y"), lua.LNumber(currentGame.readInt16(offset+14, Signed)))
		// Antal1987's dumps call this a char but has 3 unknown chars following
		L.RawSet(thisCity, lua.LString("lead_id"), lua.LNumber(currentGame.readInt8(offset+16, Signed)))
		// Unsure of these; interpreting from Antal1987's dumps
		L.RawSet(thisCity, lua.LString("improvements_maintenance"), lua.LNumber(currentGame.readInt32(offset+20, Signed)))
		L.RawSet(thisCity, lua.LString("stored_food"), lua.LNumber(currentGame.readInt32(offset+40, Signed)))
		L.RawSet(thisCity, lua.LString("stored_shields"), lua.LNumber(currentGame.readInt32(offset+44, Signed)))
		// Sample size of one, needs more testing
		numCitizens := currentGame.readInt32(offset+0x228, Signed)
		// This may be redundant since we have a table in Lua
		L.RawSet(thisCity, lua.LString("ctzn_count"), lua.LNumber(numCitizens))
		// L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
		L.RawSet(thisCity, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentGame.data[offset:offset+1024])))
		ctzn := L.NewTable()
		L.RawSet(thisCity, lua.LString("ctzn"), ctzn)
		offset += 0x22c
		for i := 0; i < numCitizens; i++ {
			lt := L.NewTable()
			ctzn.Append(lt)
			length := currentGame.readInt32(offset+4, Signed)
			// since offset is at start of "CTZN", add 8 for the [4]byte and length
			L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentGame.data[offset:offset+length+8])))
			offset += length + 8
		}

		// Only doing one city for now
		break
	}

	// L.RawSet(city, lua.LString("count"), lua.LNumber(count))
}
