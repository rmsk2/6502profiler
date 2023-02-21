package commands

import (
	"6502profiler/cpu"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type ClonedCpuProvider struct {
	cpu *cpu.CPU6502
}

func NewCloneProvider(cpu *cpu.CPU6502) (*ClonedCpuProvider, error) {
	cpu.Reset()

	return &ClonedCpuProvider{
		cpu: cpu,
	}, nil
}

func (c *ClonedCpuProvider) NewCpu() (*cpu.CPU6502, error) {
	c.cpu.Reset()
	return c.cpu, nil
}

func VerifyAllCommand(arguments []string) error {
	r := regexp.MustCompile(`^(.+)\.json$`)

	var config *cpu.Config = cpu.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verifyall", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
	preExecName := verifierFlags.String("prexec", "", "Program to run before first test")
	verboseFlag := verifierFlags.Bool("verbose", false, "Give more information")

	if err = verifierFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = cpu.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	if len(config.IoAddrConfig) != 0 {
		return fmt.Errorf("special IO addresses are incompatible with verifyall")
	}

	var cpuProv cpu.CpuProvider = config

	if *preExecName != "" {
		cpuProv, err = setupTests(config, *preExecName)
		if err != nil {
			return fmt.Errorf("unable to perform test setup: %v", err)
		}
	}

	file, err := os.Open(config.AcmeTestDir)
	if err != nil {
		return err
	}
	defer file.Close()

	names, err := file.Readdirnames(0)
	if err != nil {
		return err
	}

	testCount := 0

	for _, j := range names {
		if r.MatchString(j) {
			err := executeOneTest(j, cpuProv, config, config.AcmeTestDir, *verboseFlag)
			if err != nil {
				return err
			}

			testCount++
		}
	}

	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}
	fmt.Println()
	fmt.Printf("%d tests successfully executed\n", testCount)

	return nil
}

func setupTests(config *cpu.Config, setupPrgName string) (cpu.CpuProvider, error) {
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

	res, err := NewCloneProvider(cpu)
	if err != nil {
		return nil, fmt.Errorf("unable to perform global test setup: %v", err)
	}

	return res, nil
}

func executeOneTest(testCaseName string, cpuProv cpu.CpuProvider, asmProv cpu.AsmProvider, testDir string, verboseOutput bool) error {
	caseFileName := path.Join(testDir, testCaseName)

	if !strings.HasSuffix(caseFileName, ".json") {
		caseFileName += ".json"
	}

	testCase, err := verifier.NewTestCaseFromFile(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	cpu, err := cpuProv.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test case: %v", err)
	}
	defer func() { cpu.Mem.Close() }()

	assembler := asmProv.GetAssembler()

	if verboseOutput {
		fmt.Println("--------------------------------------------")
		fmt.Printf("Executing test case '%s'\n", testCase.Name)
		fmt.Printf("Test case file: %s\n", testCaseName)
		fmt.Printf("Test script: %s\n", testCase.TestScript)
		fmt.Printf("Test driver: %s\n", testCase.TestDriverSource)
	} else {
		fmt.Printf("Executing test case '%s' ... ", testCase.Name)
	}

	err = testCase.Execute(cpu, assembler, testDir)
	if err != nil {
		errMsg := assembler.GetErrorMessage()
		if errMsg != "" {
			fmt.Println(errMsg)
		}
		return fmt.Errorf("test case '%s' failed: %v", testCase.Name, err)
	}

	if verboseOutput {
		fmt.Printf("Clock cycles used: %d\n", cpu.NumCycles())
		fmt.Println("Test result: OK")
	} else {
		fmt.Printf("(%d clock cycles) OK\n", cpu.NumCycles())
	}

	return nil
}

func VerifyCommand(arguments []string) error {
	var config *cpu.Config = cpu.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verify", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
	testCasePath := verifierFlags.String("t", "", "Test case file")
	verboseFlag := verifierFlags.Bool("verbose", false, "Give more information")

	if err = verifierFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName != "" {
		config, err = cpu.NewConfigFromFile(*configName)
		if err != nil {
			return fmt.Errorf("error loading config: %v", err)
		}
	}

	if *testCasePath == "" {
		return fmt.Errorf("test case path has to be specified")
	}

	res := executeOneTest(*testCasePath, config, config, config.AcmeTestDir, *verboseFlag)
	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}

	return res
}
