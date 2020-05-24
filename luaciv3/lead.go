package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// Provides data from the LEAD sections which are the 32 players in a game
// Player 0 is the barbarians, player 1 is the first human player
func leadModule(L *lua.LState) {
	const numCivs = 32
	lead := L.NewTable()
	L.SetGlobal("lead", lead)
	for i := 0; i < numCivs; i++ {
		civ := L.NewTable()
		lead.Append(civ)
		leadOff, _ := currentGame.sectionOffset("LEAD", i+1)
		// TODO: I'm not sure the following is right...need to check on the offset and my relative offsets
		// queried offset is from the start of the 4-byte header, but most of my offset notes are from the end of it
		leadOff += 4
		L.RawSet(civ, lua.LString("city_count"), lua.LNumber(currentGame.readInt32(leadOff+376, Signed)))
		L.RawSet(civ, lua.LString("unit_count"), lua.LNumber(currentGame.readInt32(leadOff+368, Signed)))
		L.RawSet(civ, lua.LString("player_number"), lua.LNumber(currentGame.readInt32(leadOff+0, Signed)))
		L.RawSet(civ, lua.LString("race_id"), lua.LNumber(currentGame.readInt32(leadOff+4, Signed)))
		L.RawSet(civ, lua.LString("govt_id"), lua.LNumber(currentGame.readInt32(leadOff+132, Signed)))
		L.RawSet(civ, lua.LString("mobilization_level"), lua.LNumber(currentGame.readInt32(leadOff+136, Signed)))
		L.RawSet(civ, lua.LString("tiles_discovered"), lua.LNumber(currentGame.readInt32(leadOff+140, Signed)))
		L.RawSet(civ, lua.LString("eras_id"), lua.LNumber(currentGame.readInt32(leadOff+216, Signed)))
		L.RawSet(civ, lua.LString("research_beakers"), lua.LNumber(currentGame.readInt32(leadOff+220, Signed)))
		L.RawSet(civ, lua.LString("current_research_id"), lua.LNumber(currentGame.readInt32(leadOff+224, Signed)))
		L.RawSet(civ, lua.LString("current_research_turns"), lua.LNumber(currentGame.readInt32(leadOff+228, Signed)))
		L.RawSet(civ, lua.LString("future_techs_count"), lua.LNumber(currentGame.readInt32(leadOff+232, Signed)))
		L.RawSet(civ, lua.LString("armies_count"), lua.LNumber(currentGame.readInt32(leadOff+364, Signed)))
		L.RawSet(civ, lua.LString("military_unit_count"), lua.LNumber(currentGame.readInt32(leadOff+372, Signed)))
		atWar := L.NewTable()
		L.RawSet(civ, lua.LString("at_war"), atWar)
		for i := 0; i < numCivs; i++ {
			atWar.Append(lua.LNumber(currentGame.readInt32(leadOff+3348+(i*4), Signed)))
		}
		willTalkTo := L.NewTable()
		L.RawSet(civ, lua.LString("will_talk_to"), willTalkTo)
		for i := 0; i < numCivs; i++ {
			willTalkTo.Append(lua.LNumber(currentGame.readInt32(leadOff+2964+(i*4), Signed)))
		}
		contactWith := L.NewTable()
		L.RawSet(civ, lua.LString("contact_with"), contactWith)
		for i := 0; i < numCivs; i++ {
			contactWith.Append(lua.LNumber(currentGame.readInt32(leadOff+3732+(i*4), Signed)))
		}

	}
}
