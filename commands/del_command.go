package commands

import (
	"6502profiler/emuconfig"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
)

func DelCommand(arguments []string) error {
	delCaseFlags := flag.NewFlagSet("6502profiler delcase", flag.ContinueOnError)
	configName := delCaseFlags.String("c", "", "Config file name")
	testCasePath := delCaseFlags.String("t", "", "Test case file")

	var err error = nil
	var config *emuconfig.Config = emuconfig.DefaultConfig()

	if err := delCaseFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName == "" {
		return fmt.Errorf("a config file name has to be specified")
	}

	if *testCasePath == "" {
		return fmt.Errorf("a test case name has to be specified")
	}

	config, err = emuconfig.NewConfigFromFile(*configName)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	repo, err := config.GetCaseRepo()
	if err != nil {
		return err
	}

	err = repo.Del(*testCasePath)
	if err != nil {
		return err
	}

	return nil
}
