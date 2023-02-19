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

func VerifyAllCommand(arguments []string) error {
	r := regexp.MustCompile(`^(.+)\.json$`)

	var config *cpu.Config = cpu.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verifyall", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
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
			err := executeOneTest(j, config, *verboseFlag)
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

func executeOneTest(testCaseName string, config *cpu.Config, verboseOutput bool) error {
	caseFileName := path.Join(config.AcmeTestDir, testCaseName)

	if !strings.HasSuffix(caseFileName, ".json") {
		caseFileName += ".json"
	}

	testCase, err := verifier.NewTestCaseFromFile(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	cpu, err := config.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test case: %v", err)
	}
	defer func() { cpu.Mem.Close() }()

	assembler := config.GetAssembler()

	if verboseOutput {
		fmt.Println("--------------------------------------------")
		fmt.Printf("Executing test case '%s'\n", testCase.Name)
		fmt.Printf("Test case file: %s\n", testCaseName)
		fmt.Printf("Test script: %s\n", testCase.TestScript)
		fmt.Printf("Test driver: %s\n", testCase.TestDriverSource)
	} else {
		fmt.Printf("Executing test case '%s' ... ", testCase.Name)
	}

	err = testCase.Execute(cpu, assembler, config.AcmeTestDir)
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

	res := executeOneTest(*testCasePath, config, *verboseFlag)
	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}

	return res
}
