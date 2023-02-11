package verifier

import (
	"6502profiler/cpu"
	"encoding/hex"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type LuaCtx struct {
	cpu *cpu.CPU6502
	L   *lua.LState
}

func NewLuaCtx(cpu *cpu.CPU6502, l *lua.LState) *LuaCtx {
	return &LuaCtx{
		cpu: cpu,
		L:   l,
	}
}

func (c *LuaCtx) RegisterGlobals(L *lua.LState, loadAddress uint16, progLen uint16) error {
	L.SetGlobal("getmemory", L.NewFunction(c.GetMemory))
	L.SetGlobal("setmemory", L.NewFunction(c.SetMemory))
	L.SetGlobal("loadaddress", lua.LNumber(loadAddress))
	L.SetGlobal("proglen", lua.LNumber(progLen))

	return nil
}

func (c *LuaCtx) GetMemory(L *lua.LState) int {
	laddr := L.ToInt(1)
	memLen := L.ToInt(2)

	data, err := c.cpu.CopyFromMem(uint16(laddr), uint16(memLen))
	if err != nil {
		panic("Unable to read memory")
	}

	dataStr := hex.EncodeToString(data)
	L.Push(lua.LString(dataStr))

	return 1
}

func (c *LuaCtx) SetMemory(L *lua.LState) int {
	dataStr := L.ToString(1)
	addr := L.ToInt(2)

	data, err := hex.DecodeString(dataStr)
	if err != nil {
		panic("Unable to write memory")
	}

	err = c.cpu.CopyToMem(data, uint16(addr))
	if err != nil {
		panic("Uuable to write memory")
	}

	return 0
}

func (c *LuaCtx) callArrange() error {
	arrangeLua := lua.P{
		Fn:      c.L.GetGlobal("arrange"),
		NRet:    0,
		Protect: true,
	}

	err := c.L.CallByParam(arrangeLua)
	if err != nil {
		return fmt.Errorf("unable to call arrange function in test script: %v", err)
	}

	return nil
}

func (c *LuaCtx) callAssert() (bool, error) {
	assertLua := lua.P{
		Fn:      c.L.GetGlobal("assert"),
		NRet:    1,
		Protect: true,
	}

	err := c.L.CallByParam(assertLua)
	if err != nil {
		return false, fmt.Errorf("unable to call assert function in test script: %v", err)
	}

	ret, _ := c.L.Get(-1).(lua.LBool) // returned value
	testRes := bool(ret)
	c.L.Pop(1)

	return testRes, nil
}
