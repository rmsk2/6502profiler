package commands

import (
	"6502profiler/emuconfig"
	"6502profiler/util"
	"flag"
	"fmt"
	"os"
)

func ListCommand(arguments []string) error {
	newCaseFlags := flag.NewFlagSet("6502profiler list", flag.ContinueOnError)
	var err error = nil
	var config *emuconfig.Config = emuconfig.DefaultConfig()

	configName := newCaseFlags.String("c", "", "Config file name")

	if err := newCaseFlags.Parse(arguments); err != nil {
		os.Exit(util.ExitErrorSyntax)
	}

	if *configName == "" {
		return fmt.Errorf("a config file name has to be specified")
	}

	config, err = emuconfig.NewConfigFromFile(*configName)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	repo, err := config.GetCaseRepo()
	if err != nil {
		return err
	}

	_, err = repo.IterateTestCases(func(caseName string) error {
		testCase, err := repo.Get(caseName)
		if err != nil {
			return err
		}

		fmt.Printf("'%s'   =>   %s\n", testCase.Name, caseName)

		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to iterate over test cases: %v", err)
	}

	return nil
}
