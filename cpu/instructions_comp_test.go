package cpu

import (
	"6502profiler/memory"
	"testing"
)

func TestCmpBase(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(4096))

	cpu.Flags = 0
	cpu.cmpBase(4, 4)
	if (cpu.Flags & Flag_Z) == 0 {
		t.Fatalf("Zero flag must be set")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatalf("Negative flag must not be set")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carry flag must be set")
	}

	cpu.Flags = 0
	cpu.cmpBase(4, 5)
	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatalf("Zero flag must not be set")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatalf("Negative flag must be set")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carry flag must not be set")
	}

	cpu.Flags = 0
	cpu.cmpBase(5, 4)
	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatalf("Zero flag must not be set")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatalf("Negative flag must not be set")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carry flag must be set")
	}

	cpu.Flags = 0
	cpu.cmpBase(0x89, 0x04)
	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatalf("Zero flag must not be set")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatalf("Negative flag must be set")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carry flag must be set")
	}
}

// -------- CPY --------

func TestCPYImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 4
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// cpy #5
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xC0, 0x05, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPY immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCPYZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 4
		c.Mem.Store(0x12, 5)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// cpy $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xC4, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPY zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCPYAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 4
		c.Mem.Store(0x1214, 5)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// cpy $1214
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xCC, 0x14, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPY absolute",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- CPX --------

func TestCPXImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 4
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// cpx #4
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE0, 0x04, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPX immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCPXZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 4
		c.Mem.Store(0x12, 4)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// cpx $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE4, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPX zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCPXAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 4
		c.Mem.Store(0x1214, 4)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// cpx $1214
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xEC, 0x14, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CPX absolute",
	}

	testSingleInstructionWithCase(t, c)
}
