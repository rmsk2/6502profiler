package verifier

import (
	"6502profiler/assembler"
	"6502profiler/cpu"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

var TestDriverExtension = ".a"
var TestScriptExtension = ".lua"
var TestCaseExtension = ".json"

type SubcaseProcessor func(currentIter uint, maxIter uint)

func SetExtension(extVar *string, newVal string) {
	if !strings.HasPrefix(newVal, ".") {
		newVal = "." + newVal
	}

	*extVar = newVal
}

type TestCase struct {
	Name             string
	TestDriverSource string
	TestScript       string
}

func NewTestCase(description string, caseName string) *TestCase {
	return &TestCase{
		Name:             description,
		TestDriverSource: caseName + TestDriverExtension,
		TestScript:       caseName + TestScriptExtension,
	}
}

func NewTestCaseWithDriver(description string, caseName string, testDriverName string) *TestCase {
	return &TestCase{
		Name:             description,
		TestDriverSource: testDriverName,
		TestScript:       caseName + TestScriptExtension,
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

func (t *TestCase) Execute(cpu *cpu.CPU6502, asm assembler.Assembler, scriptPath string, subcaseProc SubcaseProcessor) error {
	var testRes bool = true
	var testMsg string
	var i uint

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

	cpu.PC = loadAdress

	err = L.DoFile(scriptToRun)
	if err != nil {
		return fmt.Errorf("unable to load test script: %v", err)
	}

	numIters, err := ctx.callNumIterations()
	if err != nil {
		numIters = 1
	}

	for i = 0; (i < numIters) && testRes; i++ {
		err = ctx.callArrange()
		if err != nil {
			return fmt.Errorf("unable to arrange test case '%s': %v", t.Name, err)
		}

		if (numIters > 1) && (subcaseProc != nil) {
			subcaseProc(i, numIters)
		}

		err = cpu.RunExt(cpu.PC, false)
		if err != nil {
			return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
		}

		testRes, testMsg, err = ctx.callAssert()
		if err != nil {
			return fmt.Errorf("unable to assert test case '%s': %v", t.Name, err)
		}
	}

	if !testRes {
		return fmt.Errorf("test failed: %s", testMsg)
	}

	return nil
}
