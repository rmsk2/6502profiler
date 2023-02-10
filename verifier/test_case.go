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

func (t *TestCase) Execute(cpu *cpu.CPU6502, assembler cpu.Assembler, testDir string) error {
	binaryToTest, err := assembler.Assemble(t.AssemblerSource)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
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
