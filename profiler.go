package main

import (
	"6502profiler/cpu"
	"6502profiler/memory"
	"fmt"
	"os"
)

func main() {
	cpu := cpu.New6502(cpu.Model6502)
	cpu.Init(memory.NewLinearMemory(16384))

	if len(os.Args) < 2 {
		fmt.Println("Usage: 6502profiler <binary to run>")
		return
	}

	res := cpu.LoadAndRun(os.Args[1])
	if res != nil {
		fmt.Printf("A problem occurred: %v\n", res)
		return
	}

	fmt.Printf("Program ran for %d clock cycles\n\n", cpu.NumCycles())
	memory.Dump(cpu.Mem, 0x0800, 0x08ff)

}
