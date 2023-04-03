package main

import (
	"6502profiler/commands"
	"6502profiler/util"
)

func main() {
	subcommParser := util.NewSubcommandParser()

	subcommParser.AddCommand("profile", commands.ProfileCommand, "Run program, record and evaulute performance data")
	subcommParser.AddCommand("verify", commands.VerifyCommand, "Run a test on an assembler program")
	subcommParser.AddCommand("verifyall", commands.VerifyAllCommand, "Run all tests")
	subcommParser.AddCommand("info", commands.InfoCommand, "Return info about program")
	subcommParser.AddCommand("newcase", commands.NewCaseCommand, "Create a new test case skeleton")
	subcommParser.AddCommand("delcase", commands.DelCommand, "Delete the files of an existing test case")
	subcommParser.AddCommand("list", commands.ListCommand, "List all test cases and their descriptions")
	subcommParser.Execute()
}
