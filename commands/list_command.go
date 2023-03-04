package commands

import (
	"6502profiler/emuconfig"
	"6502profiler/util"
	"6502profiler/verifier"
	"flag"
	"fmt"
	"os"
	"strings"
)

type listHelper struct {
	description  string
	caseFileName string
}

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

	caseList := []listHelper{}
	maxDescLen := 0

	_, err = repo.IterateTestCases(func(caseName string, testCase *verifier.TestCase) error {
		n := listHelper{
			description:  testCase.Name,
			caseFileName: caseName,
		}
		caseList = append(caseList, n)

		if len(n.description) > maxDescLen {
			maxDescLen = len(n.description)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to iterate over test cases: %v", err)
	}

	filler := strings.Repeat(" ", maxDescLen)

	for _, j := range caseList {
		fmt.Print(j.description)

		l := len(j.description)
		if l < maxDescLen {
			fmt.Print(filler[l:])
		}

		fmt.Print(" => ")
		fmt.Println(j.caseFileName)
	}

	return nil
}
