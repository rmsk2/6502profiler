package commands

import (
	"6502profiler/cpu"
	"6502profiler/emuconfig"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
	"strings"
)

type SnapshotCpuProvider struct {
	cpu *cpu.CPU6502
}

func NewSnapshotProvider(cpu *cpu.CPU6502) (*SnapshotCpuProvider, error) {
	cpu.Reset()
	cpu.Mem.TakeSnaphot()

	return &SnapshotCpuProvider{
		cpu: cpu,
	}, nil
}

func (c *SnapshotCpuProvider) NewCpu() (*cpu.CPU6502, error) {
	c.cpu.Mem.RestoreSnapshot()
	c.cpu.Reset()
	return c.cpu, nil
}

type caseExec struct {
	cpuProv     emuconfig.CpuProvider
	asmProv     emuconfig.AsmProvider
	repo        verifier.CaseRepo
	verboseFlag bool
}

func newCaseExec(c emuconfig.CpuProvider, a emuconfig.AsmProvider, repo verifier.CaseRepo, v bool) *caseExec {
	return &caseExec{
		cpuProv:     c,
		asmProv:     a,
		repo:        repo,
		verboseFlag: v,
	}
}

func (t *caseExec) loadAndExecuteCase(testCaseName string) error {
	caseFileName := testCaseName

	if !strings.HasSuffix(caseFileName, verifier.TestCaseExtension) {
		caseFileName += verifier.TestCaseExtension
	}

	testCase, err := t.repo.Get(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	return t.executeCase(testCaseName, testCase)
}

func (t *caseExec) executeCase(testCaseName string, testCase *verifier.TestCase) error {
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
			fmt.Printf("Subcase %d of %d (%d clock cycles)\n", i+1, numIter, cpu.NumCycles())
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

func VerifyAllCommand(arguments []string) error {
	var config *emuconfig.Config = emuconfig.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verifyall", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
	preExecName := verifierFlags.String("prexec", "", "Program to run before first test")
	verboseFlag := verifierFlags.Bool("verbose", false, "Give more information")

	if err = verifierFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = emuconfig.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	if len(config.IoAddrConfig) != 0 {
		return fmt.Errorf("special IO addresses are incompatible with verifyall")
	}

	repo, err := config.GetCaseRepo()
	if err != nil {
		return err
	}

	var cpuProv emuconfig.CpuProvider = config

	if *preExecName != "" {
		cpuProv, err = setupTests(config, *preExecName)
		if err != nil {
			return fmt.Errorf("unable to perform test setup: %v", err)
		}
	}

	testCount, err := repo.IterateTestCases(newCaseExec(cpuProv, config, repo, *verboseFlag).executeCase)
	if err != nil {
		return fmt.Errorf("unable to iterate test cases: %v", err)
	}

	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}
	fmt.Println()
	fmt.Printf("%d tests successfully executed\n", testCount)

	return nil
}

func setupTests(config *emuconfig.Config, setupPrgName string) (emuconfig.CpuProvider, error) {
	asm := config.GetAssembler()

	binaryName, err := asm.Assemble(setupPrgName)
	if err != nil {
		errMsg := asm.GetErrorMessage()
		if errMsg != "" {
			fmt.Println(errMsg)
		}

		return nil, fmt.Errorf("unable to setup tests: %v", err)
	}

	cpu, err := config.NewCpu()
	if err != nil {
		return nil, fmt.Errorf("unable to create cpu for test setup: %v", err)
	}

	_, _, err = cpu.LoadAndRun(binaryName)
	if err != nil {
		return nil, fmt.Errorf("unable to create cpu for test setup: %v", err)
	}

	res, err := NewSnapshotProvider(cpu)
	if err != nil {
		return nil, fmt.Errorf("unable to perform global test setup: %v", err)
	}

	return res, nil
}

func VerifyCommand(arguments []string) error {
	var config *emuconfig.Config = emuconfig.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verify", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
	testCasePath := verifierFlags.String("t", "", "Test case file")
	preExecName := verifierFlags.String("prexec", "", "Program to run before test")
	verboseFlag := verifierFlags.Bool("verbose", false, "Give more information")

	if err = verifierFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = emuconfig.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	if *testCasePath == "" {
		return fmt.Errorf("test case path has to be specified")
	}

	repo, err := config.GetCaseRepo()
	if err != nil {
		return err
	}

	var cpuProv emuconfig.CpuProvider = config

	if *preExecName != "" {
		cpuProv, err = setupTests(config, *preExecName)
		if err != nil {
			return fmt.Errorf("unable to perform test setup: %v", err)
		}
	}

	res := newCaseExec(cpuProv, config, repo, *verboseFlag).loadAndExecuteCase(*testCasePath)
	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}

	return res
}
