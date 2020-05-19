package luaciv3

import (
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// NewState is called to get a Lua environment with nbt manipulation ability// lua vm memory limit; 0 is no limit
const memoryLimitMb = 100

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
	L.SetGlobal("test", L.NewFunction(TestPassValues))

	// civ3 table
	civ3 := L.NewTable()
	L.SetGlobal("civ3", civ3)
	// if error, will get empty string, and that's fine
	path, _ := findWinCivInstall()
	L.RawSet(civ3, lua.LString("path"), lua.LString(path))

	// sav table
	sav := L.NewTable()
	L.SetGlobal("sav", sav)
	L.RawSet(sav, lua.LString("load"), L.NewFunction(SavLoad))
	L.RawSet(sav, lua.LString("dump"), L.NewFunction(SavDump))

	return nil
}

// TestPassValues is my getting familiar with calling Go from lua with values/params
//  see https://github.com/yuin/gopher-lua#calling-go-from-lua
func TestPassValues(L *lua.LState) int {
	lv := L.ToInt(1)
	lv2 := L.ToInt(2)
	L.Push(lua.LNumber(lv * lv2))
	return 1
}

// SavLoad takes a path from lua and loads it into memory
func SavLoad(L *lua.LState) int {
	path := L.ToString(1)
	err := saveGame.loadSave(path)
	// TODO: Handle errors
	if err != nil {
		panic(err)
	}
	sav := L.GetGlobal("sav")
	if savTable, ok := sav.(*lua.LTable); ok {
		L.RawSet(savTable, lua.LString("path"), lua.LString(saveGame.path))
		L.RawSet(savTable, lua.LString("name"), lua.LString(saveGame.fileName()))
	}
	return 0
}

// TODO: parameters
// SavDump returns a hex dump to lua
func SavDump(L *lua.LState) int {
	dump := hex.Dump(saveGame.data[:256])
	L.Push(lua.LString(dump))
	return 1
}
