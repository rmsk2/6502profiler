package verifier

import (
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

func (t *TestCase) Execute(cpu *cpu.CPU6502, assembler cpu.Assembler, testDir string) error {
	binaryToTest, err := assembler.Assemble(t.TestDriverSource)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	loadAdress, progLen, err := cpu.Load(binaryToTest)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	scriptPath := path.Join(testDir, t.TestScript)

	L := lua.NewState()
	defer L.Close()

	ctx := NewLuaCtx(cpu, L)

	err = L.DoFile(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to load test script: %v", err)
	}

	err = ctx.RegisterGlobals(L, loadAdress, progLen)
	if err != nil {
		return fmt.Errorf("unable to register Lua functions: %v", err)
	}

	fmt.Printf("Executing test case '%s' ... ", t.Name)

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

	fmt.Println("OK")

	return nil
}
