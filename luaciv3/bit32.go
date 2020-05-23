package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// This is a partial reimplementation of Lua 5.2's bit32 module which is not available in gopher-lua made with the 5.1 specs
// http://www.lua.org/manual/5.2/manual.html#6.7

func bit32Module(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("bit32", lt)
	L.RawSet(lt, lua.LString("band"), L.NewFunction(bit32Band))
}

// variadic lua-callable function to bitwise-and all operands
func bit32Band(L *lua.LState) int {
	var ops []int
	var numOps int
	var res int
	// Since variadic, will convert to string to check for nil first
	for {
		op := L.ToString(numOps + 1)
		if op == "" {
			break
		}
		numOps++
		ops = append(ops, L.ToInt(numOps))
	}
	if numOps > 0 {
		res = ops[0]
		for i := 1; i < numOps; i++ {
			res &= ops[i]
		}
	}
	L.Push(lua.LNumber(res))
	return 1
}
