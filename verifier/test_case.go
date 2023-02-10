package verifier

import (
	"6502profiler/cpu"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path"

	lua "github.com/yuin/gopher-lua"
)

type TestCase struct {
	Name            string
	AssemblerSource string
	TestScript      string
}

func NewTestCaseFromFile(fileName string) (*TestCase, error) {
	var res *TestCase = &TestCase{}

	testCaseData, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to load testcase file %s: %v", fileName, err)
	}

	err = json.Unmarshal(testCaseData, res)
	if err != nil {
		return nil, fmt.Errorf("unable to load testcase file %s: %v", fileName, err)
	}

	return res, nil
}

func (t *TestCase) Save(fileName string) error {
	data, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to save testcase file %s: %v", fileName, err)
	}

	err = os.WriteFile(fileName, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to save config testcase file %s: %v", fileName, err)
	}

	return nil
}

func callArrange(L *lua.LState) error {
	arrangeLua := lua.P{
		Fn:      L.GetGlobal("arrange"),
		NRet:    0,
		Protect: true,
	}

	err := L.CallByParam(arrangeLua)
	if err != nil {
		return fmt.Errorf("unable to call arrange function in test script: %v", err)
	}

	return nil
}

func callAssert(L *lua.LState) (bool, error) {
	assertLua := lua.P{
		Fn:      L.GetGlobal("assert"),
		NRet:    1,
		Protect: true,
	}

	err := L.CallByParam(assertLua)
	if err != nil {
		return false, fmt.Errorf("unable to call assert function in test script: %v", err)
	}

	ret, _ := L.Get(-1).(lua.LBool) // returned value
	testRes := bool(ret)
	L.Pop(1)

	return testRes, nil
}

type LuaCtx struct {
	cpu *cpu.CPU6502
}

func (c *LuaCtx) RegisterFunctions(L *lua.LState) error {
	L.SetGlobal("getmemory", L.NewFunction(c.GetMemory))
	L.SetGlobal("setmemory", L.NewFunction(c.SetMemory))

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

func (t *TestCase) Execute(cpu *cpu.CPU6502, assembler cpu.Assembler, testDir string) error {
	binaryToTest, err := assembler.Assemble(t.AssemblerSource)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	ctx := &LuaCtx{
		cpu: cpu,
	}

	loadAdress, _, err := cpu.Load(binaryToTest)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	scriptPath := path.Join(testDir, t.TestScript)

	L := lua.NewState()
	defer L.Close()

	err = L.DoFile(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to load test script: %v", err)
	}

	err = ctx.RegisterFunctions(L)
	if err != nil {
		return fmt.Errorf("unable to register Lua functions: %v", err)
	}

	fmt.Printf("Executing test case '%s' ... ", t.Name)

	err = callArrange(L)
	if err != nil {
		return fmt.Errorf("unable to arrange test case '%s': %v", t.Name, err)
	}

	err = cpu.Run(loadAdress)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	testRes, err := callAssert(L)
	if err != nil {
		return fmt.Errorf("unable to assert test case '%s': %v", t.Name, err)
	}

	if !testRes {
		fmt.Println("FAIL")
		return fmt.Errorf("test failed")
	}

	fmt.Println("OK")

	return nil
}
