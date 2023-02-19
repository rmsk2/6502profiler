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

func NewTestCase(testName string, fileNamePrefix string) *TestCase {
	return &TestCase{
		Name:             testName,
		TestDriverSource: fileNamePrefix + ".a",
		TestScript:       fileNamePrefix + ".lua",
	}
}

func NewTestCaseWithDriver(testName string, fileNamePrefix string, testDriverName string) *TestCase {
	return &TestCase{
		Name:             testName,
		TestDriverSource: testDriverName,
		TestScript:       fileNamePrefix + ".lua",
	}
}

func (t *TestCase) WriteSekeleton(fileNamePrefix string, testDir string, createDriver bool) error {
	scriptPath := path.Join(testDir, t.TestScript)
	testDriverPath := path.Join(testDir, t.TestDriverSource)
	jsonPath := path.Join(testDir, fileNamePrefix+".json")

	_, err := os.Stat(jsonPath)
	if err == nil {
		return fmt.Errorf("json file '%s' already exists", jsonPath)
	}

	_, err = os.Stat(scriptPath)
	if err == nil {
		return fmt.Errorf("script file '%s' already exists", scriptPath)
	}

	if createDriver {
		_, err = os.Stat(testDriverPath)
		if err == nil {
			return fmt.Errorf("test driver file '%s' already exists", testDriverPath)
		}
	}

	data, err := json.MarshalIndent(t, "", "    ")
	if err != nil {
		return fmt.Errorf("unable to save testcase file %s: %v", jsonPath, err)
	}

	err = os.WriteFile(jsonPath, data, 0600)
	if err != nil {
		return fmt.Errorf("unable to save config testcase file %s: %v", jsonPath, err)
	}

	f, err := os.Create(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to create lua script '%s'", scriptPath)
	}
	defer func() { f.Close() }()

	if createDriver {
		f2, err := os.Create(testDriverPath)
		if err != nil {
			return fmt.Errorf("unable to create test driver '%s'", testDriverPath)
		}
		defer func() { f2.Close() }()
	}

	return nil
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

	ctx := NewLuaCtx(cpu, testDir, L)

	err = ctx.RegisterGlobals(L, loadAdress, progLen)
	if err != nil {
		return fmt.Errorf("unable to register Lua functions: %v", err)
	}

	err = L.DoFile(scriptPath)
	if err != nil {
		return fmt.Errorf("unable to load test script: %v", err)
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

	fmt.Printf("(%d clock cycles) ... OK \n", cpu.NumCycles())

	return nil
}
