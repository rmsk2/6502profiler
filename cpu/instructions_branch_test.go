package cpu

import "testing"

// -------- BPL --------

func TestBPLSkip(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 0
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
		c.Flags = 0xFF
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
		c.Flags = 255
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

func TestBMIUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 0
		c.Flags |= Flag_N
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bmi .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x30, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BMI",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBEQUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 0
		c.Flags |= Flag_Z
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    beq .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0xF0, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BEQ",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBNEUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 255
		c.Flags &= (^Flag_Z)
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bne .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0xd0, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BNE",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBCCUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 255
		c.Flags &= (^Flag_C)
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bcc .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x90, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BCC",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBCSUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 0
		c.Flags |= Flag_C
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bcs .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0xB0, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BCS",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBVCUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 255
		c.Flags &= (^Flag_V)
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bvc .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x50, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BVC",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBVSUp(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 0
		c.Flags |= Flag_V
		c.PC = UnitProgStart + 1
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	//.up
	//    brk
	//    bvs .up
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x00, 0x70, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BVS",
	}

	testSingleInstructionWithCase(t, c)
}
