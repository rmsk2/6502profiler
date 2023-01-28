package cpu

import (
	"6502profiler/memory"
	"testing"
)

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

// -------- add/subtract --------

func TestAddBin(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBin(1, 1)

	if r != 2 {
		t.Fatal("Basic binary add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for binary add")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for binary add")
	}

	if (cpu.Flags & Flag_V) != 0 {
		t.Fatal("Overflow flag not correct for binary add")
	}
}

func TestAddBinRollOver(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBin(0xFF, 1)

	if r != 0 {
		t.Fatal("Basic binary add does not work")
	}

	if (cpu.Flags & Flag_Z) == 0 {
		t.Fatal("Zero flag not correct for binary add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for binary add")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for binary add")
	}

	if (cpu.Flags & Flag_V) != 0 {
		t.Fatal("Overflow flag not correct for binary add")
	}
}

func TestAddBinOverflowSet(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBin(208, 144)

	if r != 96 {
		t.Fatal("Basic binary add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for binary add")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for binary add")
	}

	if (cpu.Flags & Flag_V) == 0 {
		t.Fatal("Overflow flag not correct for binary add")
	}
}

func TestAddBinOverflowSetNoCarry(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBin(80, 80)

	if r != 160 {
		t.Fatal("Basic binary add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary add")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Negative flag not correct for binary add")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for binary add")
	}

	if (cpu.Flags & Flag_V) == 0 {
		t.Fatal("Overflow flag not correct for binary add")
	}
}

func TestSUBBin(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBin(2, 1)

	if r != 1 {
		t.Fatal("Basic binary sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for binary sub")
	}

	if (cpu.Flags & Flag_V) != 0 {
		t.Fatal("Overflow flag not correct for binary sub")
	}
}

func TestSUBBinRollover(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBin(0, 1)

	if r != 0xFF {
		t.Fatal("Basic binary sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Negative flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for binary sub")
	}

	if (cpu.Flags & Flag_V) != 0 {
		t.Fatal("Overflow flag not correct for binary sub")
	}
}

func TestSUBBinNoBorrowButOverflow(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBin(0xD0, 0x70)

	if r != 96 {
		t.Fatal("Basic binary sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for binary sub")
	}

	if (cpu.Flags & Flag_V) == 0 {
		t.Fatal("Overflow flag not correct for binary sub")
	}
}

func TestSUBBinBorrowAndOverflow(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBin(0x50, 0xB0)

	if r != 160 {
		t.Fatal("Basic binary sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Negative flag not correct for binary sub")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for binary sub")
	}

	if (cpu.Flags & Flag_V) == 0 {
		t.Fatal("Overflow flag not correct for binary sub")
	}
}

func TestAddBaseBCD(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBcd6502(0x19, 0x22)

	if r != 0x41 {
		t.Fatal("Basic BCD add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for BCD add")
	}
}

func TestAddBaseBCDRollover(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.addBaseBcd6502(0x99, 0x10)

	if r != 0x09 {
		t.Fatal("Basic BCD add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for BCD add")
	}
}

func TestAddBaseBCDZero(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.addBaseBcd6502(0x99, 0x00)

	if r != 0x00 {
		t.Fatal("Basic BCD add does not work")
	}

	if (cpu.Flags & Flag_Z) == 0 {
		t.Fatal("Zero flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for BCD add")
	}
}

func TestAddBaseBCDRolloverCarrySet(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.addBaseBcd6502(0x99, 0x10)

	if r != 0x10 {
		t.Fatal("Basic BCD add does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD add")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for BCD add")
	}
}

func TestSubBaseBCD(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBcd(0x22, 0x19)

	if r != 0x03 {
		t.Fatal("Basic BCD sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Carryflag not correct for BCD sub")
	}
}

func TestSubBaseBCDNeg(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	r := cpu.subBaseBcd(0x19, 0x22)

	if r != 0x97 {
		t.Fatal("Basic BCD sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Negative flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for BCD sub")
	}
}

func TestSubBaseBCDNegWithBorrow(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.subBaseBcd(0x19, 0x22)

	if r != 0x96 {
		t.Fatal("Basic BCD sub does not work")
	}

	if (cpu.Flags & Flag_Z) != 0 {
		t.Fatal("Zero flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Negative flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for BCD sub")
	}
}

func TestSubBaseBCDMaxNeg(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	cpu.Flags = 0xFF
	cpu.Flags &= (^Flag_C)
	r := cpu.subBaseBcd(0x00, 0x99)

	if r != 0x00 {
		t.Fatal("Basic BCD sub does not work")
	}

	if (cpu.Flags & Flag_Z) == 0 {
		t.Fatal("Zero flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_N) != 0 {
		t.Fatal("Negative flag not correct for BCD sub")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Carryflag not correct for BCD sub")
	}
}

// -------- ADC --------

func TestADCImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x01
		c.Flags &= (^Flag_C)
	}

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

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// adc #$41
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x69, 0x41, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCImmediateBCD(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags &= (^Flag_C)
		c.Flags |= Flag_D
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x30 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// adc #$11
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x69, 0x11, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC immediate BCD",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags &= (^Flag_C)
		c.Flags |= Flag_D
		c.X = 0x02
		c.Mem.Store(0x0012, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x30 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// adc $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x35, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC Zero page X BCD",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- SBC --------

func TestSBCImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x43
		c.Flags |= Flag_C
	}

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

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc #$01
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE9, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCImmediateBCD(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x40
		c.Flags |= Flag_C
		c.Flags |= Flag_D
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x39 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// sbc #$01
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE9, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC immediate BCD",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags |= Flag_C
		c.Mem.Store(0x0012, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x08 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// sbc $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE5, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCBCDZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags &= (^Flag_C)
		c.Flags |= Flag_D
		c.X = 0x02
		c.Mem.Store(0x0012, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x30 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) != 0 {
			return false
		}

		if (c.Flags & Flag_C) != 0 {
			return false
		}

		return true
	}

	// adc $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x35, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC Zero page X BCD",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCBCDZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags |= Flag_C
		c.Flags |= Flag_D
		c.X = 0x02
		c.Mem.Store(0x0012, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x08 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// sbc $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xF5, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC Zero page X BCD",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Flags &= (^Flag_C)
		c.Mem.Store(0x1000, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x01 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// adc $1000
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x6D, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Flags |= Flag_C
		c.Mem.Store(0x1000, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0xDF {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc $1000
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xED, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.X = 0x02
		c.Flags &= (^Flag_C)
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x01 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// adc $1000, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x7D, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.X = 0x02
		c.Flags |= Flag_C
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0xDF {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc $1000, x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xFD, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Y = 0x02
		c.Flags &= (^Flag_C)
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x01 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// adc $1000, y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x79, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC absolute Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Y = 0x02
		c.Flags |= Flag_C
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0xDF {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc $1000, y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xF9, 0x00, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC absolute Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCIndirectIdxY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Y = 0x02
		c.Flags &= (^Flag_C)
		c.Mem.Store(0x0012, 0x00)
		c.Mem.Store(0x0013, 0x10)
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x01 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// adc ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x71, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC indirect index Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCIndirectIdxY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.Y = 0x02
		c.Flags |= Flag_C
		c.Mem.Store(0x0012, 0x00)
		c.Mem.Store(0x0013, 0x10)
		c.Mem.Store(0x1002, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0xDF {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xF1, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC indirect index Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestADCIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.X = 0x02
		c.Flags &= (^Flag_C)
		c.Mem.Store(0x0012, 0x00)
		c.Mem.Store(0x0013, 0x10)
		c.Mem.Store(0x1000, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0x01 {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
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

	// adc ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x61, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC index X indirect",
	}

	testSingleInstructionWithCase(t, c)
}

func TestSBCIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF0
		c.X = 0x02
		c.Flags |= Flag_C
		c.Mem.Store(0x0012, 0x00)
		c.Mem.Store(0x0013, 0x10)
		c.Mem.Store(0x1000, 0x11)
	}

	verifier := func(c *CPU6502) bool {
		if c.A != 0xDF {
			return false
		}

		if (c.Flags & Flag_Z) != 0 {
			return false
		}

		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_C) == 0 {
			return false
		}

		return true
	}

	// sbc ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE1, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "SBC index X indirect",
	}

	testSingleInstructionWithCase(t, c)
}
