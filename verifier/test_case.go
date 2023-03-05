package verifier

import (
	"6502profiler/assembler"
	"6502profiler/cpu"
	"encoding/json"
	"fmt"
	"os"
	"path"

	lua "github.com/yuin/gopher-lua"
)

type TestCase struct {
	Name             string
	TestDriverSource string
	TestScript       string
}

func NewTestCase(description string, caseName string) *TestCase {
	return &TestCase{
		Name:             description,
		TestDriverSource: caseName + ".a",
		TestScript:       caseName + ".lua",
	}
}

func NewTestCaseWithDriver(description string, caseName string, testDriverName string) *TestCase {
	return &TestCase{
		Name:             description,
		TestDriverSource: testDriverName,
		TestScript:       caseName + ".lua",
	}
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

func (t *TestCase) Execute(cpu *cpu.CPU6502, asm assembler.Assembler, scriptPath string) error {
	binaryToTest, err := asm.Assemble(t.TestDriverSource)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	loadAdress, progLen, err := cpu.Load(binaryToTest)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	scriptToRun := path.Join(scriptPath, t.TestScript)

	L := lua.NewState()
	defer L.Close()

	ctx := NewLuaCtx(cpu, scriptPath, L)

	err = ctx.RegisterGlobals(L, loadAdress, progLen)
	if err != nil {
		return fmt.Errorf("unable to register Lua functions: %v", err)
	}

	err = L.DoFile(scriptToRun)
	if err != nil {
		return fmt.Errorf("unable to load test script: %v", err)
	}

	err = ctx.callArrange()
	if err != nil {
		return fmt.Errorf("unable to arrange test case '%s': %v", t.Name, err)
	}

	err = cpu.Run(loadAdress)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	testRes, testMsg, err := ctx.callAssert()
	if err != nil {
		return fmt.Errorf("unable to assert test case '%s': %v", t.Name, err)
	}

	if !testRes {
		return fmt.Errorf("test failed: %s", testMsg)
	}

	return nil
}
