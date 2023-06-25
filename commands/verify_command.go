package commands

import (
	"6502profiler/caseexec"
	"6502profiler/emuconfig"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
)

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
		cpuProv, err = config.SetupTests(*preExecName)
		if err != nil {
			return fmt.Errorf("unable to perform test setup: %v", err)
		}
	}

	testCount, err := repo.IterateTestCases(caseexec.NewCaseExec(cpuProv, config, repo, *verboseFlag).ExecuteCase)
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
		cpuProv, err = config.SetupTests(*preExecName)
		if err != nil {
			return fmt.Errorf("unable to perform test setup: %v", err)
		}
	}

	res := caseexec.NewCaseExec(cpuProv, config, repo, *verboseFlag).LoadAndExecuteCase(*testCasePath)
	if *verboseFlag {
		fmt.Println("--------------------------------------------")
	}

	return res
}
