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
		// offset is right, but not sure about length; 20 is what Antal1987's dumps say
		L.RawSet(thisCity, lua.LString("name"), lua.LString(civString(currentGame.data[offset+0x188:offset+0x188+20])))
		// L.RawSet(lt, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentBic.data[off:off+length])))
		L.RawSet(thisCity, lua.LString("dump"), lua.LString("\n"+hex.Dump(currentGame.data[offset:offset+512])))
		// Only doing one city for now
		break
	}

	// L.RawSet(city, lua.LString("count"), lua.LNumber(count))
}
