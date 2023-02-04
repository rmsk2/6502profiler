package cpu

import (
	"6502profiler/memory"
	"testing"
)

// A VerifyFunc returns true if the test is OK
type VerifyFunc func(c *CPU6502) bool

// A PrepareFunc prepares the machine for a unit test
type PrepareFunc func(c *CPU6502)

// Start address for all unit testing programs
const UnitProgStart = 0x0800

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

func testSingleInstructionWithArrange(model CpuModel, testProg []byte, arranger PrepareFunc, verifier VerifyFunc) (bool, error) {
	cpu := New6502(model)
	cpu.Init(memory.NewLinearMemory(8192))

	err := cpu.CopyBinary(testProg, 0x800)
	if err != nil {
		return false, err
	}

	cpu.PC = UnitProgStart

	if arranger != nil {
		arranger(cpu)
	}

	err = cpu.Run(cpu.PC)
	if err != nil {
		return false, err
	}

	res := verifier(cpu)

	return res, err
}

type InstructionTestCase struct {
	model           CpuModel
	testProg        []byte
	arranger        PrepareFunc
	verifier        VerifyFunc
	instructionName string
}

func testSingleInstructionWithCase(t *testing.T, c InstructionTestCase) {
	res, err := testSingleInstructionWithArrange(c.model, c.testProg, c.arranger, c.verifier)

	if res == false {
		t.Fatalf("%s could not be verified", c.instructionName)
	}
	if err != nil {
		t.Fatalf("%s does not work: %v", c.instructionName, err)
	}
}
