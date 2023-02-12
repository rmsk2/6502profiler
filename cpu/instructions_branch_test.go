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

func TestBRAUp(t *testing.T) {
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
	//    bra .up
	//    brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x00, 0x80, 0xfd, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BRA up",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBRADown(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Flags = 255
		c.Flags &= (^Flag_N)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+4
	}

	//    bra .down
	//    brk
	//.down
	//    brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x80, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BRA down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR0Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0xFE)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr0 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x0f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR0 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR0NotDown(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0xFF) // Make sure branch does not happen
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+4
	}

	// 	bbr0 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x0f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR0 down (not branching)",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR0Up(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0xFE)
		c.PC = 0x0801
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	// .back
	// 		brk
	// 		lda #$FE
	// 		sta $40
	// 		bbr0 $40, .back
	// 		brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x00, 0xa9, 0xfe, 0x85, 0x40, 0x0f, 0x40, 0xf8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR0 up",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR1Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x02
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr1 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x1f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR1 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR2Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x04
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr2 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x2f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR2 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR3Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x08
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr3 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x3f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR3 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR4Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x10
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr4 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x4f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR4 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR5Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x20
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr5 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x5f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR5 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR6Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x40
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr6 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x6f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR6 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBR7Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		var val uint8 = 0xFF ^ 0x80
		c.Mem.Store(0x0040, val)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbr $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x7f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBR7 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS0Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x01)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs0 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x8f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS0 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS0NotDown(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x00) // Make sure branch does not happen
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+4
	}

	// 	bbs0 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x8f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS0 down (not branching)",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS0Up(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x01)
		c.PC = 0x0801
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+1
	}

	// .back
	// 		brk
	// 		lda #$01
	// 		sta $40
	// 		bbs0 $40, .back
	// 		brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x00, 0xa9, 0x01, 0x85, 0x40, 0x8f, 0x40, 0xf8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS0 up",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS1Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x02)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs1 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x9f, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS1 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS2Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x04)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs2 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xaf, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS2 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS3Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x08)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs3 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xbf, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS3 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS4Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x10)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs4 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xcf, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS4 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS5Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x20)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs5 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xdf, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS5 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS6Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x40)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs6 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xef, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS6 down",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBBS7Down(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0040, 0x80)
	}

	verifier := func(c *CPU6502) bool {
		return c.PC == UnitProgStart+5
	}

	// 	bbs7 $40, skip
	// 	brk
	//skip
	// 	brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0xFf, 0x40, 0x01, 0x00, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BBS7 down",
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
