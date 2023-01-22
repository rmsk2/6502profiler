package cpu

import (
	"6502profiler/memory"
	"testing"
)

// A VerifyFunc returns true if the test is OK
type VerifyFunc func(c *CPU6502) bool

func TestStackOperations(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(3072))

	if cpu.SP != 0xFF {
		t.Fatal("Stack pointer not at start value")
	}

	cpu.push(0x42)
	if cpu.SP != 0xFE {
		t.Fatal("Stack pointer incorrect after push")
	}

	if cpu.Mem.Load(0x1FF) != 0x42 {
		t.Fatal("Value was not pushed correctly")
	}

	res := cpu.pop()
	if res != 0x42 {
		t.Fatal("Pop does not work")
	}

	if cpu.SP != 0xFF {
		t.Fatal("Stack pointer incorrect after pop")
	}
}

func testSingleInstruction(model CpuModel, testProg []byte, verifier VerifyFunc) (bool, error) {
	cpu := New6502(model)
	cpu.Init(memory.NewLinearMemory(8192))

	err := cpu.CopyProg(testProg, 0x800)
	if err != nil {
		return false, err
	}

	err = cpu.Run(0x800)
	if err != nil {
		return false, err
	}

	res := verifier(cpu)

	return res, err
}
