package cpu

import "testing"

// -------- LDA --------

func TestLDAImmediate(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.A != 0x42 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// lda #$42
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA9, 0x42, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDA immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAImmediate0(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.A != 0x00 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// lda #00
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA9, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDA immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAImmediateNeg(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.A != 0x81 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		return true
	}

	// lda #$81
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA9, 0x81, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDA immediate",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsolute(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda $0804
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xAD, 0x04, 0x08, 0x00, 0x72},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDA absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x12, 0x72)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA5, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x20, 0x72)
		c.X = 0x0E
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda $12, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xB5, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA zero page index X",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 4
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda $0800, x
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xBD, 0x00, 0x08, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA absolute with X index",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDA implementations
// => no extra test
func TestLDAAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 4
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda $0800, y
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xB9, 0x00, 0x08, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA absolute with Y index",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAIndirectIdxY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 3
		c.Mem.Store(0x0012, 0x00)
		c.Mem.Store(0x0013, 0x08)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda ($12),y
	// brk
	// !byte $72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xb1, 0x12, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA indirect with Y index",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 3
		c.Mem.Store(0x0012, 0x03)
		c.Mem.Store(0x0013, 0x08)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda ($12)
	// brk
	// !byte $72
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xb2, 0x12, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA indirect",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDAIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 2
		c.Mem.Store(0x0012, 0x03)
		c.Mem.Store(0x0013, 0x08)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x72
	}

	// lda ($10, x)
	// brk
	// !byte $72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA1, 0x10, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDA X index indirect",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- LDX --------

func TestLDXImmediate(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.X != 0x42 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// ldx #$42
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA2, 0x42, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDX immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDXImmediate0(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.X != 0x00 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// ldx #00
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA2, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDX immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDXImmediateNeg(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.X != 0x81 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		return true
	}

	// ldx #$81
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA2, 0x81, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDX immediate",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDX implementations
// => no extra test
func TestLDXAbsolute(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.X == 0x72
	}

	// ldx $0804
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xAE, 0x04, 0x08, 0x00, 0x72},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDX absolute",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDX implementations
// => no extra test
func TestLDXAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 4
	}

	verifier := func(c *CPU6502) bool {
		return c.X == 0x72
	}

	// ldx $0800, y
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xBE, 0x00, 0x08, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDX absolute with Y index",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDXZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x72)
	}

	verifier := func(c *CPU6502) bool {
		return c.X == 0x72
	}

	// ldx $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA6, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDX Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDXZeroPageY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(230, 0x72)
		c.Y = 10
	}

	verifier := func(c *CPU6502) bool {
		return c.X == 0x72
	}

	// ldx $DC, y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xB6, 220, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDX Zero page with Y",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- LDY --------

func TestLDYImmediate(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.Y != 0x42 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// ldy #$42
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA0, 0x42, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDY immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDYImmediate0(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.Y != 0x00 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		return true
	}

	// ldy #00
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA0, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDY immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDYImmediateNeg(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		if c.Y != 0x81 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		return true
	}

	// ldy #$81
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA0, 0x81, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDY immediate",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDY implementations
// => no extra test
func TestLDYAbsolute(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.Y == 0x72
	}

	// ldy $0804
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xAC, 0x04, 0x08, 0x00, 0x72},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "LDY absolute",
	}

	testSingleInstructionWithCase(t, c)
}

// Code to set N and Z flags is the same in all LDY implementations
// => no extra test
func TestLDYAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 4
	}

	verifier := func(c *CPU6502) bool {
		return c.Y == 0x72
	}

	// ldy $0800, x
	// brk
	// !byte 0x72
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xBC, 0x00, 0x08, 0x00, 0x72},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDY absolute with X index",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDYZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x72)
	}

	verifier := func(c *CPU6502) bool {
		return c.Y == 0x72
	}

	// ldy $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA4, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDY Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLDYZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(230, 0x72)
		c.X = 10
	}

	verifier := func(c *CPU6502) bool {
		return c.Y == 0x72
	}

	// ldy $DC, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xB4, 220, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LDY Zero page with X",
	}

	testSingleInstructionWithCase(t, c)
}
