package util

import (
	"fmt"
	"os"
	"sort"
)

const (
	ExitError       int = 43
	ExitErrorSyntax int = 42
	ExitOk          int = 0
)

// CommandFunc is a type for functions that implement a subcommand
type CommandFunc func(arguments []string) error

type commandInfo struct {
	command  CommandFunc
	usageMsg string
}

// SubCommParser provides the basis for a flag parser that treats the first argument as a command
type SubCommParser struct {
	knownCommands map[string]*commandInfo
}

// NewSubcommandParser returns an initialized subcommand flag parser
func NewSubcommandParser() *SubCommParser {
	return &SubCommParser{
		knownCommands: map[string]*commandInfo{},
	}
}

// AddCommand adds a subcommand to the parser
func (s *SubCommParser) AddCommand(name string, f CommandFunc, usage string) {
	s.knownCommands[name] = &commandInfo{
		command:  f,
		usageMsg: usage,
	}
}

// Execute parses the command line and calls the appropriate command function
func (s *SubCommParser) Execute() {
	switch commandLineLength := len(os.Args); {
	case commandLineLength < 2:
		s.PrintDefaults()
		os.Exit(ExitErrorSyntax)
	case commandLineLength == 2:
		if (os.Args[1] == "-h") || (os.Args[1] == "-?") {
			s.PrintDefaults()
			os.Exit(ExitErrorSyntax)
		}
		fallthrough
	default:
		subCommand, ok := s.knownCommands[os.Args[1]]
		if !ok {
			s.PrintDefaults()
			os.Exit(ExitErrorSyntax)
		}

		err := subCommand.command(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			os.Exit(ExitErrorSyntax)
		}
	}
}

// PrintDefaults prints a description of the availbale commands
func (s *SubCommParser) PrintDefaults() {
	keys := make([]string, 0, len(s.knownCommands))
	for i := range s.knownCommands {
		keys = append(keys, i)
	}

	sort.Strings(keys)

	fmt.Println("The following commands are available: ")
	for _, j := range keys {
		fmt.Printf("     %s: %s\n", j, s.knownCommands[j].usageMsg)
	}
}
