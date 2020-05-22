package luaciv3

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

// NewState is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 200

// NewState creates and returns a lua state with LuaCiv3 functions
func NewState() *lua.LState {
	L := lua.NewState()
	// Set memory limit of lua instance (just a safety measure)
	if memoryLimitMb > 0 {
		L.SetMx(memoryLimitMb)
	}
	if err := LuaCiv3(L); err != nil {
		// TODO: Better error handling, although not expecting setup runtime errors
		fmt.Printf("\n********************\n%s\n\n", err.Error())
	}
	return L
}

// LuaCiv3 injects functions into a gopher-lua state
// TODO: Should I eliminate error return? lua.NewState() doesn't return error
func LuaCiv3(L *lua.LState) error {
	// test function
	L.SetGlobal("test", L.NewFunction(testPassValues))

	// install_path
	// if error, will get empty string, and that's fine
	path, _ := findWinCivInstall()
	L.SetGlobal("install_path", lua.LString(path))

	// sav table
	sav := L.NewTable()
	L.SetGlobal("sav", sav)
	L.RawSet(sav, lua.LString("load"), L.NewFunction(savLoad))
	L.RawSet(sav, lua.LString("dump"), L.NewFunction(savDump))

	// bic table
	bic := L.NewTable()
	L.SetGlobal("bic", bic)
	L.RawSet(bic, lua.LString("load_default"), L.NewFunction(bicLoadDefault))
	L.RawSet(bic, lua.LString("dump"), L.NewFunction(bicDump))

	L.SetGlobal("get_savs", L.NewFunction(getSavs))

	return nil
}

// testPassValues is my getting familiar with calling Go from lua with values/params
//  see https://github.com/yuin/gopher-lua#calling-go-from-lua
func testPassValues(L *lua.LState) int {
	lv := L.ToInt(1)
	lv2 := L.ToInt(2)
	L.Push(lua.LNumber(lv * lv2))
	return 1
}

// TODO: Sav- and Bic- loads and dumps are similar; perhaps move to struct method?
// savLoad takes a path from lua and loads it into memory
func savLoad(L *lua.LState) int {
	path := L.ToString(1)
	err := saveGame.loadSave(path)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}
	L.SetGlobal("save_path", lua.LString(saveGame.path))
	L.SetGlobal("save_name", lua.LString(saveGame.fileName()))

	civ3Module(L)
	tileModule(L)
	leadModule(L)
	prtoModule(L)
	raceModule(L)
	unitModule(L)
	gameModule(L)
	techModule(L)
	erasModule(L)
	govtModule(L)
	wsizModule(L)

	return 0
}

// TODO: parameters
// savDump returns a hex dump to lua
func savDump(L *lua.LState) int {
	dump := hex.Dump(saveGame.data[:30])
	L.Push(lua.LString(dump))
	return 1
}

// bicLoadDefault takes a path from lua and loads it into memory
func bicLoadDefault(L *lua.LState) int {
	path := L.ToString(1)
	// Try to fetch install_path if no path provided
	if path == "" {
		installPath := L.GetGlobal("install_path")
		if iPathString, ok := installPath.(lua.LString); ok {
			path = string(iPathString) + "/conquests.biq"
		}
	}
	err := defaultBic.loadSave(path)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}
	L.SetGlobal("bic_path", lua.LString(defaultBic.path))
	L.SetGlobal("bic_name", lua.LString(defaultBic.fileName()))
	return 0
}

// TODO: parameters
// bicDump returns a hex dump to lua
func bicDump(L *lua.LState) int {
	dump := hex.Dump(defaultBic.data[:256])
	L.Push(lua.LString(dump))
	return 1
}

// getSavs takes a table as input, each value is a directory for which to return all .SAV file paths
func getSavs(L *lua.LState) int {
	dirs := L.ToTable(1)
	savs := L.NewTable()
	dirs.ForEach(func(_ lua.LValue, v lua.LValue) {
		if dir, ok := v.(lua.LString); ok {
			if files, err := ioutil.ReadDir(string(dir)); err == nil {
				for _, f := range files {
					if (!f.IsDir()) && strings.ToLower(filepath.Ext(f.Name())) == ".sav" {
						savs.Append(lua.LString(string(dir) + "/" + f.Name()))
					}
				}
			} else {
				fmt.Println("whoops")
			}
		}
	})
	L.Push(savs)
	return 1
}
