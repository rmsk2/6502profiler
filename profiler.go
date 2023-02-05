package main

import (
	"6502profiler/acmeassembler"
	"6502profiler/cpu"
	"6502profiler/memory"
	"fmt"
	"os"
	"strconv"
)

const (
	ExitError int = 43
	ExitOk    int = 0
)

func main() {
	cpu := cpu.New6502(cpu.Model6502)
	// mem := memory.NewMemWrapper(memory.NewLinearMemory(16384), 0x2D00)
	// picProc := memory.NewPicProcessor(320, 200)
	// mem.AddSpecialWriteAddress(0x2DDD, picProc.SetPoint)
	mem := memory.NewLinearMemory(65536)
	var p float64

	cpu.Init(mem)

	if !((len(os.Args) == 2) || (len(os.Args) == 5)) {
		fmt.Println("Usage: 6502profiler <binary to run> [<label file> <file to save runtime info> <cutoff percent>]")
		os.Exit(ExitError)
	}

	if len(os.Args) >= 5 {
		percent, err := strconv.ParseUint(os.Args[4], 10, 7)

		if err != nil {
			fmt.Printf("Can not parse percentage: %v\n", err)
			os.Exit(ExitError)
		}

		if percent > 100 {
			fmt.Printf("Percentage too large: %d\n", percent)
			os.Exit(ExitError)
		}

		p = float64(percent) / 100.0
	}

	res := cpu.LoadAndRun(os.Args[1])
	if res != nil {
		fmt.Printf("A problem occurred: %v\n", res)
		os.Exit(ExitError)
	}

	fmt.Printf("Program ran for %d clock cycles\n", cpu.NumCycles())

	if len(os.Args) >= 5 {
		labels, err := acmeassembler.ParseLabelFile(os.Args[2])
		if err != nil {
			fmt.Printf("A problem occurred: %v\n", err)
			os.Exit(ExitError)
		}

		ctOff := func(m memory.Memory, start, end uint16) uint64 {
			//return memory.CutOffAbsoluteValue(m, start, end, p)
			return memory.CutOffMedian(m, start, end, p)
		}

		memory.DumpStatistics(mem, os.Args[3], labels, 2048, 2048+2048, ctOff)
	}

	// picProc.Save("apfel.png")

	os.Exit(ExitOk)
}
