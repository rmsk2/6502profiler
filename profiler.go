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
	mem := memory.NewPicWrapper(memory.NewLinearMemory(16384), 320, 200)
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
	memory.Dump(cpu.Mem, 0x0800, 0x08ff)

	mem.Save("apfel.png")

	os.Exit(ExitOk)
}
