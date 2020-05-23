package luaciv3

import (
	lua "github.com/yuin/gopher-lua"
)

// This is a partial reimplementation of Lua 5.2's bit32 module which is not available in gopher-lua made with the 5.1 specs
// http://www.lua.org/manual/5.2/manual.html#6.7

// There are issues with this module regarding int size, but for the use cases in this package it should suffice
//   More info: bit32 is supposed to work on unsigned 32-bit integers, but Lua
//   numbers are float64s, and this converts them to a plantform-dependent signed
//   int, so there are potential issues with all that conversion.

func bit32Module(L *lua.LState) {
	lt := L.NewTable()
	L.SetGlobal("bit32", lt)
	L.RawSet(lt, lua.LString("band"), L.NewFunction(bit32band))
	L.RawSet(lt, lua.LString("rshift"), L.NewFunction(bit32rshift))
	L.RawSet(lt, lua.LString("bnot"), L.NewFunction(bit32bnot))
	L.RawSet(lt, lua.LString("bor"), L.NewFunction(bit32bor))
}

// variadic lua-callable function to bitwise-and all operands
func bit32band(L *lua.LState) int {
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

// Shift first operand to the right 2nd operand bits
func bit32bnot(L *lua.LState) int {
	x := L.ToInt(1)
	res := ^x
	res &= 0xffffffff
	L.Push(lua.LNumber(res))
	return 1
}

// variadic lua-callable function to bitwise-or all operands
func bit32bor(L *lua.LState) int {
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
			res |= ops[i]
		}
	}
	L.Push(lua.LNumber(res))
	return 1
}

// Shift first operand to the left 2nd operand bits
func bit32lshift(L *lua.LState) int {
	x := L.ToInt(1)
	disp := L.ToInt(2)
	res := x << disp
	L.Push(lua.LNumber(res))
	return 1
}

// Shift first operand to the right 2nd operand bits
func bit32rshift(L *lua.LState) int {
	x := L.ToInt(1)
	disp := L.ToInt(2)
	res := x >> disp
	L.Push(lua.LNumber(res))
	return 1
}

// arshift, btest, bxor, extract, replace, lrotate, and rrotate not implemented
