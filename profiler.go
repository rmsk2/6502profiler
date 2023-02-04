package main

import (
	"6502profiler/acmeassembler"
	"6502profiler/cpu"
	"6502profiler/memory"
	"fmt"
	"os"
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
	mem := memory.NewLinearMemory(16384)
	cpu.Init(mem)

	if len(os.Args) < 3 {
		fmt.Println("Usage: 6502profiler <binary to run> <label file>")
		os.Exit(ExitError)
	}

	labels, err := acmeassembler.ParseLabelFile(os.Args[2])
	if err != nil {
		fmt.Printf("A problem occurred: %v\n", err)
		os.Exit(ExitError)
	}

	res := cpu.LoadAndRun(os.Args[1])
	if res != nil {
		fmt.Printf("A problem occurred: %v\n", res)
		os.Exit(ExitError)
	}

	fmt.Printf("Program ran for %d clock cycles\n", cpu.NumCycles())
	ctOff := func(m memory.Memory, start, end uint16) uint64 {
		//return memory.CutOffAbsoluteValue(m, start, end, 0.1)
		return memory.CutOffMedian(m, start, end, 0.1)
	}

	memory.DumpStatistics(mem, "access_data.txt", labels, 2048, 2048+2048, ctOff)

	// picProc.Save("apfel.png")

	os.Exit(ExitOk)
}
