package caseexec

import (
	"6502profiler/emuconfig"
	"6502profiler/verifier"
	"fmt"
	"strings"
)

type CaseExec struct {
	cpuProv     emuconfig.CpuProvider
	asmProv     emuconfig.AsmProvider
	repo        verifier.CaseRepo
	verboseFlag bool
}

func NewCaseExec(c emuconfig.CpuProvider, a emuconfig.AsmProvider, repo verifier.CaseRepo, v bool) *CaseExec {
	return &CaseExec{
		cpuProv:     c,
		asmProv:     a,
		repo:        repo,
		verboseFlag: v,
	}
}

func (t *CaseExec) LoadAndExecuteCase(testCaseName string) error {
	caseFileName := testCaseName

	if !strings.HasSuffix(caseFileName, verifier.TestCaseExtension) {
		caseFileName += verifier.TestCaseExtension
	}

	testCase, err := t.repo.Get(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	return t.ExecuteCase(testCaseName, testCase)
}

func (t *CaseExec) ExecuteCase(testCaseName string, testCase *verifier.TestCase) error {
	cpu, err := t.cpuProv.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test case: %v", err)
	}
	defer func() { cpu.Mem.Close() }()

	assembler := t.asmProv.GetAssembler()
	var subcaseProc verifier.SubcaseProcessor = nil

	if t.verboseFlag {
		fmt.Println("--------------------------------------------")
		fmt.Printf("Executing test case '%s'\n", testCase.Name)
		fmt.Printf("Test case file: %s\n", testCaseName)
		fmt.Printf("Test script: %s\n", testCase.TestScript)
		fmt.Printf("Test driver: %s\n", testCase.TestDriverSource)

		subcaseProc = func(i uint, numIter uint) {
			fmt.Printf("Executing subcase %d of %d (%d clock cycles already used)\n", i+1, numIter, cpu.NumCycles())
		}
	} else {
		fmt.Printf("Executing test case '%s' ... ", testCase.Name)
	}

	err = testCase.Execute(cpu, assembler, t.repo.GetScriptPath(), subcaseProc)
	if err != nil {
		errMsg := assembler.GetErrorMessage()
		if errMsg != "" {
			fmt.Println(errMsg)
		}
		return fmt.Errorf("test case '%s' failed: %v", testCase.Name, err)
	}

	if t.verboseFlag {
		fmt.Printf("Clock cycles used: %d\n", cpu.NumCycles())
		fmt.Println("Test result: OK")
	} else {
		fmt.Printf("(%d clock cycles) OK\n", cpu.NumCycles())
	}

	return nil

}
