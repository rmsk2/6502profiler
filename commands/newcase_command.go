package commands

import (
	"6502profiler/cpu"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
)

func NewCaseCommand(arguments []string) error {
	newCaseFlags := flag.NewFlagSet("6502profiler newcase", flag.ContinueOnError)
	var err error = nil
	var config *cpu.Config = cpu.DefaultConfig()

	configName := newCaseFlags.String("c", "", "Config file name")
	prefixName := newCaseFlags.String("p", "", "File name prefix")
	testName := newCaseFlags.String("n", "", "Test name")

	if err := newCaseFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName == "" {
		return fmt.Errorf("a config file name has to be specified")
	}

	if *prefixName == "" {
		return fmt.Errorf("a prefix has to be specified")
	}

	if *testName == "" {
		return fmt.Errorf("a test name has to be specified")
	}

	config, err = cpu.NewConfigFromFile(*configName)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	testCase := verifier.NewTestCase(*testName, *prefixName)

	err = testCase.WriteSekeleton(*prefixName, config.AcmeTestDir)
	if err != nil {
		return fmt.Errorf("error writing skeleton: %v", err)
	}

	return nil
}
