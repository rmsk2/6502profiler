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
	L.SetGlobal("get_memory", L.NewFunction(c.GetMemory))
	L.SetGlobal("set_memory", L.NewFunction(c.SetMemory))
	L.SetGlobal("load_address", lua.LNumber(loadAddress))
	L.SetGlobal("prog_len", lua.LNumber(progLen))
	L.SetGlobal("get_flags", L.NewFunction(c.GetFlagsLua))
	L.SetGlobal("set_flags", L.NewFunction(c.SetFlagsLua))
	L.SetGlobal("get_sp", L.NewFunction(c.GetSP))
	L.SetGlobal("get_pc", L.NewFunction(c.GetPC))
	L.SetGlobal("get_accu", L.NewFunction(c.GetAccu))
	L.SetGlobal("get_xreg", L.NewFunction(c.GetX))
	L.SetGlobal("get_yreg", L.NewFunction(c.GetY))
	L.SetGlobal("set_accu", L.NewFunction(c.SetAccu))
	L.SetGlobal("set_xreg", L.NewFunction(c.SetX))
	L.SetGlobal("set_yreg", L.NewFunction(c.SetY))

	return nil
}

func (c *LuaCtx) GetSP(L *lua.LState) int {
	return c.GetRegister(L, &c.cpu.SP)
}

func (c *LuaCtx) GetAccu(L *lua.LState) int {
	return c.GetRegister(L, &c.cpu.A)
}

func (c *LuaCtx) GetX(L *lua.LState) int {
	return c.GetRegister(L, &c.cpu.X)
}

func (c *LuaCtx) GetY(L *lua.LState) int {
	return c.GetRegister(L, &c.cpu.Y)
}

func (c *LuaCtx) SetAccu(L *lua.LState) int {
	return c.SetRegister(L, &c.cpu.A)
}

func (c *LuaCtx) SetX(L *lua.LState) int {
	return c.SetRegister(L, &c.cpu.X)
}

func (c *LuaCtx) SetY(L *lua.LState) int {
	return c.SetRegister(L, &c.cpu.Y)
}

func (c *LuaCtx) GetPC(L *lua.LState) int {
	L.Push(lua.LNumber(c.cpu.PC))

	return 1
}

func (c *LuaCtx) GetRegister(L *lua.LState, reg *uint8) int {
	L.Push(lua.LNumber(*reg))

	return 1
}

func (c *LuaCtx) SetRegister(L *lua.LState, reg *uint8) int {
	newValue := uint8(L.ToInt(1))

	*reg = newValue

	return 0
}

func (c *LuaCtx) GetFlags() string {
	res := []byte{}

	if (c.cpu.Flags & cpu.Flag_N) != 0 {
		res = append(res, 'N')
	} else {
		res = append(res, '-')
	}

	if (c.cpu.Flags & cpu.Flag_V) != 0 {
		res = append(res, 'V')
	} else {
		res = append(res, '-')
	}

	res = append(res, '-')

	if (c.cpu.Flags & cpu.Flag_B) != 0 {
		res = append(res, 'B')
	} else {
		res = append(res, '-')
	}

	if (c.cpu.Flags & cpu.Flag_D) != 0 {
		res = append(res, 'D')
	} else {
		res = append(res, '-')
	}

	if (c.cpu.Flags & cpu.Flag_I) != 0 {
		res = append(res, 'I')
	} else {
		res = append(res, '-')
	}

	if (c.cpu.Flags & cpu.Flag_Z) != 0 {
		res = append(res, 'Z')
	} else {
		res = append(res, '-')
	}

	if (c.cpu.Flags & cpu.Flag_C) != 0 {
		res = append(res, 'C')
	} else {
		res = append(res, '-')
	}

	return string(res)
}

func (c *LuaCtx) SetFlags(flags string) {
	var res uint8 = 0

	if len(flags) > 8 {
		panic("flag value is too large")
	}

	for _, j := range flags {
		switch j {
		case 'N':
			res |= cpu.Flag_N
		case 'V':
			res |= cpu.Flag_V
		case 'B':
			res |= cpu.Flag_B
		case 'D':
			res |= cpu.Flag_D
		case 'I':
			res |= cpu.Flag_I
		case 'Z':
			res |= cpu.Flag_Z
		case 'C':
			res |= cpu.Flag_C
		}

	}

	c.cpu.Flags = res
}

func (c *LuaCtx) GetFlagsLua(L *lua.LState) int {
	L.Push(lua.LString(c.GetFlags()))

	return 1
}

func (c *LuaCtx) SetFlagsLua(L *lua.LState) int {
	flagsStr := L.ToString(1)

	c.SetFlags(flagsStr)

	return 0
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