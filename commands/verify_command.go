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
			err := executeOneTest(j, config)
			if err != nil {
				return err
			}

			testCount++
		}
	}

	fmt.Println()
	fmt.Printf("%d tests successfully executed\n", testCount)

	return nil
}

func executeOneTest(testCaseName string, config *cpu.Config) error {
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
	err = testCase.Execute(cpu, assembler, config.AcmeTestDir)
	if err != nil {
		errMsg := assembler.GetErrorMessage()
		if errMsg != "" {
			fmt.Println(errMsg)
		}
		return fmt.Errorf("test case '%s' failed: %v", testCase.Name, err)
	}

	return nil
}

func VerifyCommand(arguments []string) error {
	var config *cpu.Config = cpu.DefaultConfig()
	var err error
	verifierFlags := flag.NewFlagSet("6502profiler verify", flag.ContinueOnError)
	configName := verifierFlags.String("c", "", "Config file name")
	testCasePath := verifierFlags.String("t", "", "Test case file")

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

	return executeOneTest(*testCasePath, config)
}
