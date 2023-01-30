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

func TestSTAAbsoluteX(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x10A8, 0x44)
		c.A = 0x52
		c.X = 0xA8
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x10A8) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta $1000, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x9D, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA absolute with X index",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTAZeroPageX(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0043, 0x44)
		c.X = 0x10
		c.A = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0043) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta $33, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x95, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTAIndirectY(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1057, 0x44)
		c.Mem.Store(0x0012, 0x50)
		c.Mem.Store(0x0013, 0x10)
		c.Y = 0x7
		c.A = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x1057) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x91, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA indirect Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTAXIndirect(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1057, 0x44)
		c.Mem.Store(0x0012, 0x57)
		c.Mem.Store(0x0013, 0x10)
		c.X = 0x02
		c.A = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x1057) == 0x52) && (flagsBefore == c.Flags)
	}

	// sta ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x81, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STA X indirect",
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

// -------- STX --------

func TestSTXZeroPage(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0033, 0x44)
		c.X = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0033) == 0x52) && (flagsBefore == c.Flags)
	}

	// stx $33
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x86, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STX zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTXZeroPageY(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0043, 0x44)
		c.Y = 0x10
		c.X = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0043) == 0x52) && (flagsBefore == c.Flags)
	}

	// stx $33, y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x96, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STX zero page Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTXAbsolute(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1000, 0x44)
		c.X = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x1000) == 0x52) && (flagsBefore == c.Flags)
	}

	// stx $1000
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x8E, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STX absolute",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- STY --------

func TestSTYZeroPageX(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0043, 0x44)
		c.X = 0x10
		c.Y = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0043) == 0x52) && (flagsBefore == c.Flags)
	}

	// sty $33, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x94, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STY zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTYAbsolute(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1000, 0x44)
		c.Y = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x1000) == 0x52) && (flagsBefore == c.Flags)
	}

	// sty $1000
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x8C, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STY absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSTYZeroPage(t *testing.T) {
	var flagsBefore uint8

	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0033, 0x44)
		c.Y = 0x52
		flagsBefore = c.Flags
	}

	verifier := func(c *CPU6502) bool {
		return (c.Mem.Load(0x0033) == 0x52) && (flagsBefore == c.Flags)
	}

	// sty $33
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x84, 0x33, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "STY zero page",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- Flag stuff --------

func TestCLC(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= Flag_C
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_C) == 0
	}

	// clc
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x18, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CLC",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCLD(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= Flag_D
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_D) == 0
	}

	// cld
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xD8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CLD",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCLV(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= Flag_V
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_V) == 0
	}

	// clv
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xB8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CLV",
	}

	testSingleInstructionWithCase(t, c)
}

func TestCLI(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= Flag_I
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_I) == 0
	}

	// cli
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x58, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "CLI",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSEC(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags &= (^Flag_C)
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_C) != 0
	}

	// sec
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x38, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SEC",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSED(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags &= (^Flag_D)
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_D) != 0
	}

	// sed
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xF8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SED",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSEI(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags &= (^Flag_I)
	}

	verifier := func(c *CPU6502) bool {
		return (c.Flags & Flag_I) != 0
	}

	// sei
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x78, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SEI",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- PHP --------

func TestPHP(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= (^Flag_I)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x01FF) == c.Flags
	}

	// php
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x08, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PHP",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- PLP --------

func TestPLP(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= (^Flag_I)
		c.push(c.Flags)
		c.Flags = 0
	}

	verifier := func(c *CPU6502) bool {
		return c.Flags == ^Flag_I
	}

	// plp
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x28, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "PLP",
	}

	testSingleInstructionWithCase(t, c)
}
