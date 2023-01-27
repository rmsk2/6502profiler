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

// -------- jsr/rts --------

func TestJSR(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return (c.A == 0x42) && (c.X == 6)
	}

	//    lda #5
	//    jsr .overwrite
	//    ldx #6
	//    brk
	//.overwrite
	//    lda #0x42
	//    rts
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xa9, 0x05, 0x20, 0x08, 0x08, 0xa2, 0x06, 0x00, 0xa9, 0x42, 0x60},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "JSR/RTS",
	}

	testSingleInstructionWithCase(t, c)
}

func TestJMP(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.PC == 0x0807
	}

	//    jmp testLocation
	//    brk
	//    brk
	//    brk
	//testLocation
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x4c, 0x06, 0x08, 0x00, 0x00, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "JMP",
	}

	testSingleInstructionWithCase(t, c)
}

func TestJMPIndirect6502(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.PC == 0x0808
	}

	//    jmp (jmpAddr)
	//    brk
	//jmpAddr
	//!byte <jmpTarget, >jmpTarget
	//    brk
	//jmpTarget
	//    brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x6c, 0x04, 0x08, 0x00, 0x07, 0x08, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "JMP indirect 6502",
	}

	testSingleInstructionWithCase(t, c)
}

func TestJMPIndirect65C02(t *testing.T) {
	verifier := func(c *CPU6502) bool {
		return c.PC == 0x0808
	}

	//    jmp (jmpAddr)
	//    brk
	//jmpAddr
	//!byte <jmpTarget, >jmpTarget
	//    brk
	//jmpTarget
	//    brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x6c, 0x04, 0x08, 0x00, 0x07, 0x08, 0x00, 0x00},
		arranger:        nil,
		verifier:        verifier,
		instructionName: "JMP indirect 6502",
	}

	testSingleInstructionWithCase(t, c)
}
