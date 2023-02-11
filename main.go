package main

import (
	"6502profiler/commands"
	"6502profiler/util"
)

func main() {
	subcommParser := util.NewSubcommandParser()

	subcommParser.AddCommand("profile", commands.ProfileCommand, "Record data about program executions")
	subcommParser.AddCommand("verify", commands.VerifyCommand, "Run a test on an assembler program")
	subcommParser.AddCommand("verifyall", commands.VerifyAllCommand, "Run all tests")
	subcommParser.Execute()
}
