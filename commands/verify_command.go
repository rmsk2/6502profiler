package commands

import (
	"6502profiler/cpu"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
	"path"
)

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

	caseFileName := path.Join(config.AcmeTestDir, *testCasePath)

	testCase, err := verifier.NewTestCaseFromFile(caseFileName)
	if err != nil {
		return fmt.Errorf("unable to load test case file: %v", err)
	}

	cpu, err := config.NewCpu()
	if err != nil {
		return fmt.Errorf("unable to create cpu for test case: %v", err)
	}

	err = testCase.Execute(cpu, config.GetAssembler(), config.AcmeSrcDir)
	if err != nil {
		return fmt.Errorf("test case '%s' failed: %v", testCase.Name, err)
	}

	return nil
}
