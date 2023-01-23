package cpu

import "testing"

// -------- BPL --------

func TestBPLSkip(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags |= Flag_N
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+4
	}

	//.up
	//    brk
	//    bpl .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x10, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BPL skip",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBPLUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags &= (^Flag_N)
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bpl .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x10, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BPL up",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBPLDown(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags &= (^Flag_N)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+4
	}

	//    bpl .down
	//    brk
	//.down
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x10, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BPL down",
	}

	testSingleInstructionWithCase(t, c)
}
