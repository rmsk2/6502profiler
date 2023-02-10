package main

import (
	"6502profiler/commands"
	"6502profiler/util"
)

func main() {
	subcommParser := util.NewSubcommandParser()

	subcommParser.AddCommand("profile", commands.ProfileCommand, "Record data about program executions")
	subcommParser.AddCommand("verify", commands.VerifyCommand, "Run tests on assembler programs ")
	subcommParser.Execute()
}
