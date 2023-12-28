package luabridge

import (
	"6502profiler/cpu"
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type TrapProcessor struct {
	Ctx *LuaCtx
}

func NewTrapProcessor(l *lua.LState, scriptToRun string, cpu *cpu.CPU6502, loadAddress uint16, progLen uint16, id string) (*TrapProcessor, error) {
	res := &TrapProcessor{
		Ctx: NewLuaCtx(cpu, "", l),
	}

	res.Ctx.SetIdent(id)

	err := res.Ctx.RegisterGlobals(l, loadAddress, progLen)
	if err != nil {
		return nil, fmt.Errorf("unable to register Lua functions: %v", err)
	}

	err = l.DoFile(scriptToRun)
	if err != nil {
		return nil, fmt.Errorf("unable to load script: %v", err)
	}

	return res, nil
}

func (t *TrapProcessor) Write(trapCode uint8) {
	err := t.Ctx.CallTrap(trapCode)
	if err != nil {
		panic(fmt.Sprintf("unable to call trap function in Lua code: %v", err))
	}
}
