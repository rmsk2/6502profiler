package main

import (
	"6502profiler/acmeassembler"
	"6502profiler/cpu"
	"6502profiler/memory"
	"6502profiler/profiler"
	"6502profiler/util"
	"fmt"
	"os"
	"strconv"
)

func ProfileCommand(arguments []string) error {
	cpu := cpu.New6502(cpu.Model6502)
	// mem := memory.NewMemWrapper(memory.NewLinearMemory(16384), 0x2D00)
	// picProc := memory.NewPicProcessor(320, 200)
	// mem.AddSpecialWriteAddress(0x2DDD, picProc.SetPoint)
	mem := memory.NewLinearMemory(65536)
	var p float64

	cpu.Init(mem)

	if !((len(arguments) == 1) || (len(arguments) == 4)) {
		fmt.Println("Usage: 6502profiler <binary to run> [<label file> <file to save runtime info> <cutoff percent>]")
		os.Exit(util.ExitError)
	}

	if len(arguments) >= 4 {
		percent, err := strconv.ParseUint(arguments[3], 10, 7)

		if err != nil {
			fmt.Printf("Can not parse percentage: %v\n", err)
			os.Exit(util.ExitError)
		}

		if percent > 100 {
			fmt.Printf("Percentage too large: %d\n", percent)
			os.Exit(util.ExitError)
		}

		p = float64(percent) / 100.0
	}

	res := cpu.LoadAndRun(arguments[0])
	if res != nil {
		fmt.Printf("A problem occurred: %v\n", res)
		os.Exit(util.ExitError)
	}

	fmt.Printf("Program ran for %d clock cycles\n", cpu.NumCycles())

	if len(os.Args) >= 5 {
		labels, err := acmeassembler.ParseLabelFile(arguments[1])
		if err != nil {
			fmt.Printf("A problem occurred: %v\n", err)
			os.Exit(util.ExitError)
		}

		ctOff := func(m memory.Memory, start, end uint16) uint64 {
			//return memory.CutOffAbsoluteValue(m, start, end, p)
			return profiler.CutOffMedian(m, start, end, p)
		}

		profiler.DumpStatistics(mem, arguments[2], labels, 2048, 2048+2048, ctOff)
	}

	// picProc.Save("apfel.png")

	os.Exit(util.ExitOk)

	return nil
}

func main() {
	subcommParser := util.NewSubcommandParser()

	subcommParser.AddCommand("profile", ProfileCommand, "Generate data about program executions")
	subcommParser.Execute()
}
