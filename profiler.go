package main

import (
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
	//mem := memory.NewMemWrapper(memory.NewLinearMemory(16384), 0x2D00)
	//picProc := memory.NewPicProcessor(320, 200)
	//mem.AddSpecialWriteAddress(0x2DDD, picProc.SetPoint)
	mem := memory.NewLinearMemory(16384)
	cpu.Init(mem)

	if len(os.Args) < 2 {
		fmt.Println("Usage: 6502profiler <binary to run>")
		os.Exit(ExitError)
	}

	res := cpu.LoadAndRun(os.Args[1])
	if res != nil {
		fmt.Printf("A problem occurred: %v\n", res)
		os.Exit(ExitError)
	}

	fmt.Printf("Program ran for %d clock cycles\n\n", cpu.NumCycles())
	memory.DumpStatistics(mem, "access_data.txt", 2048, 2048+3072)

	//picProc.Save("apfel.png")

	os.Exit(ExitOk)
}
