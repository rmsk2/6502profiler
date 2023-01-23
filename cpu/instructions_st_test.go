package cpu

import "testing"

// -------- STA --------

func TestSTAAbsolute(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1000, 0x44)
		c.A = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x1000) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta $1000
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x8D, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTAZeroPage(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0033, 0x44)
		c.A = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0033) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta $33
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x85, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTAAbsoluteY(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x10A8, 0x44)
		c.A = 0x52
		c.Y = 0xA8
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x10A8) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta $1000, y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x99, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA absolute with Y index",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- PHA --------

func TestPHA(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.A = 0x52
		c.SP = 0xFF
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		if flagsBefore != c.Flags {
			return false
		}

		if c.Mem.Load(0x1FF) != 0x52 {
			return false
		}

		if c.SP != 0xFE {
			return false
		}

		return true
	}

	// pha
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x48, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PHA",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- PLA --------

func TestPLA(t *testing.T) {

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0120, 0x17)
		c.SP = 0x1F
	}

	verifier := func(c *CPU6502) bool {
		return (c.A == 0x17) && (c.SP == 0x20)
	}

	// pla
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x68, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PLA",
	}

	testSingleInstructionWithCase(t, c)
}

func TestPLA0(t *testing.T) {

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0100, 0x85)
		c.SP = 0xFF
	}

	verifier := func(c *CPU6502) bool {
		if (c.A != 0x85) || (c.SP != 0x00) {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		return (c.Flags & Flag_N) != 0
	}

	// pla
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x68, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PLA",
	}

	testSingleInstructionWithCase(t, c)
}

func TestPLANeg(t *testing.T) {

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0101, 0x00)
		c.SP = 0x00
	}

	verifier := func(c *CPU6502) bool {
		if (c.A != 0x00) || (c.SP != 0x01) {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return (c.Flags & Flag_Z) != 0
	}

	// pla
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x68, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PLA",
	}

	testSingleInstructionWithCase(t, c)
}
