package luaciv3

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// LuaCiv3 injects functions into a gopher-lua state
func LuaCiv3(L *lua.LState) error {
	fmt.Println("luaciv3 doesn't do anything yet")
	return nil
}
