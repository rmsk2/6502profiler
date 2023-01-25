package cpu

import "testing"

// -------- DEY --------

func TestDEY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 0x43
	}

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

	// dey
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x88, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEY",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDEYNeg(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 0x00
	}

	verifier := func(c *CPU6502) bool {
		if c.Y != 0xFF {
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

	// dey
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x88, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEY",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDEY0(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 0x01
		c.Flags |= Flag_N // Check whether negative flag is reset
	}

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

	// dey
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x88, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEY",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- INY --------

// Don't test setting of zero and negative flag. INY uses the same
// code (nzFlags()) as the other instructions which have already been
// tested.
func TestINY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 0x43
	}

	verifier := func(c *CPU6502) bool {
		if c.Y != 0x44 {
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

	// iny
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xC8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INY",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDEX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 0x43
	}

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

	// dex
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xCA, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEX",
	}

	testSingleInstructionWithCase(t, c)
}

func TestINX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 0x43
	}

	verifier := func(c *CPU6502) bool {
		if c.X != 0x44 {
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

	// inx
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INX",
	}

	testSingleInstructionWithCase(t, c)
}
