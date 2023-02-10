package verifier

import (
	"6502profiler/cpu"
	"encoding/json"
	"fmt"
	"os"
	"path"
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
	fmt.Println(scriptPath)

	// Initialize scripting environment

	// Load script file

	fmt.Printf("Executing test case '%s' ... ", t.Name)
	// Call arrange() function in script

	err = cpu.Run(loadAdress)
	if err != nil {
		return fmt.Errorf("unable to execute test case '%s': %v", t.Name, err)
	}

	// Call assert() function in script

	fmt.Println("OK")

	return nil
}
