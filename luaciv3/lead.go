package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

const numCivs = 32

// Was unsure what to implement next in the project, got a little inspiration from
//  https://forums.civfanatics.com/threads/when-do-tier-2-barbarian-units-spawn.657845/
//  So I'll try to count cities, barb units and types, and barb camps
//  Eventually this will be refactored to other modules, but for now we'll call it suede after OP
func leadModule(L *lua.LState) {
	lead := L.NewTable()
	L.SetGlobal("lead", lead)
	civs := L.NewTable()
	L.RawSet(lead, lua.LString("civs"), civs)
	for i := 0; i < numCivs; i++ {
		civ := L.NewTable()
		civs.Append(civ)
		leadOff, _ := currentGame.sectionOffset("LEAD", i+1)
		// queried offset is from the star of the 4-byte header, but most of my offset notes are from the end of it
		leadOff += 4
		L.RawSet(civ, lua.LString("city_count"), lua.LNumber(currentGame.readInt32(leadOff+376, Signed)))
		L.RawSet(civ, lua.LString("unit_count"), lua.LNumber(currentGame.readInt32(leadOff+368, Signed)))
		L.RawSet(civ, lua.LString("player_number"), lua.LNumber(currentGame.readInt32(leadOff+0, Signed)))
		L.RawSet(civ, lua.LString("race_id"), lua.LNumber(currentGame.readInt32(leadOff+4, Signed)))
		L.RawSet(civ, lua.LString("government_type"), lua.LNumber(currentGame.readInt32(leadOff+132, Signed)))
		L.RawSet(civ, lua.LString("mobilization_level"), lua.LNumber(currentGame.readInt32(leadOff+136, Signed)))
		L.RawSet(civ, lua.LString("tiles_discovered"), lua.LNumber(currentGame.readInt32(leadOff+140, Signed)))
		L.RawSet(civ, lua.LString("era"), lua.LNumber(currentGame.readInt32(leadOff+252, Signed)))
		L.RawSet(civ, lua.LString("research_beakers"), lua.LNumber(currentGame.readInt32(leadOff+220, Signed)))
		L.RawSet(civ, lua.LString("current_researchId"), lua.LNumber(currentGame.readInt32(leadOff+224, Signed)))
		L.RawSet(civ, lua.LString("current_researchTurns"), lua.LNumber(currentGame.readInt32(leadOff+228, Signed)))
		L.RawSet(civ, lua.LString("future_techs_count"), lua.LNumber(currentGame.readInt32(leadOff+232, Signed)))
		L.RawSet(civ, lua.LString("armies_count"), lua.LNumber(currentGame.readInt32(leadOff+364, Signed)))
		L.RawSet(civ, lua.LString("military_unit_count"), lua.LNumber(currentGame.readInt32(leadOff+372, Signed)))

	}
}
