package luaciv3

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// NewState creates and returns a lua state with LuaCiv3 functions
func NewState() *lua.LState {
	L := lua.NewState()
	_ = LuaCiv3(L)
	return L
}

// LuaCiv3 injects functions into a gopher-lua state
// TODO: Should I eliminate error return? lua.NewState() doesn't return error
func LuaCiv3(L *lua.LState) error {
	fmt.Println("luaciv3 doesn't do anything yet")
	return nil
}
