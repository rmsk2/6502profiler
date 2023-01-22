package cpu

import (
	"6502profiler/memory"
	"testing"
)

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
