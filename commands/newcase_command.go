package commands

import (
	"6502profiler/emuconfig"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
)

func NewCaseCommand(arguments []string) error {
	newCaseFlags := flag.NewFlagSet("6502profiler newcase", flag.ContinueOnError)
	var err error = nil
	var config *emuconfig.Config = emuconfig.DefaultConfig()

	configName := newCaseFlags.String("c", "", "Config file name")
	testCaseName := newCaseFlags.String("p", "", "Test case file name")
	testDescription := newCaseFlags.String("d", "", "Test description")
	testDriverName := newCaseFlags.String("t", "", "Full name of test driver file in test dir (optional)")
	var testCase *verifier.TestCase

	if err := newCaseFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName == "" {
		return fmt.Errorf("a config file name has to be specified")
	}

	if *testCaseName == "" {
		return fmt.Errorf("a test case file name has to be specified")
	}

	if *testDescription == "" {
		return fmt.Errorf("a test description has to be specified")
	}

	config, err = emuconfig.NewConfigFromFile(*configName)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	repo, err := config.GetCaseRepo()
	if err != nil {
		return err
	}

	if *testDriverName == "" {
		testCase = verifier.NewTestCase(*testDescription, *testCaseName)
	} else {
		testCase = verifier.NewTestCaseWithDriver(*testDescription, *testCaseName, *testDriverName)
	}

	err = repo.New(*testCaseName, testCase, *testDriverName == "")
	if err != nil {
		return fmt.Errorf("error creating new test case: %v", err)
	}

	return nil
}
