package luaciv3

import (
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Provides data from the RACE sectionsof the BIC
func raceModule(L *lua.LState) {
	// var name string
	race := L.NewTable()
	L.SetGlobal("race", race)
	raceOff, _ := currentBic.sectionOffset("RACE", 1)
	listSection(raceOff, func(off int) {
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
	})
	fmt.Println(hex.Dump(currentBic.data[raceOff : raceOff+256]))

	/*
		        race {
		            leaderName: string(offset:0, maxLength: 32)
		            leaderTitle: string(offset:32, maxLength: 24)
		            adjective:  string(offset:88, maxLength: 40)
		            civName: string(offset:128, maxLength: 40)
		            objectNoun: string(offset:168, maxLength: 40)
		        }

			    tradeRace: race {
		        civName: string(offset:128, maxLength: 40)
			}


				prtoOff, _ := currentBic.sectionOffset("RACE", 1)
				numPrto := currentBic.readInt32(prtoOff, Signed)
				off := prtoOff + 4
				fmt.Println(numPrto)
				for i := 0; i < numPrto; i++ {
					lt := L.NewTable()
					prto.Append(lt)
					prtoLen = currentBic.readInt32(off, Signed)
					// skip over the length
					off += 4
					name, err := civString(currentBic.data[off+4 : off+4+32])
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
	*/
}
