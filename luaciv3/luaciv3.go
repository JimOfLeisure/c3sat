package luaciv3

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// Blua injects functions into a gopher-lua state
func Blua(L *lua.LState) error {
	fmt.Println("luaciv3 doesn't do anything yet")
	return nil
}
