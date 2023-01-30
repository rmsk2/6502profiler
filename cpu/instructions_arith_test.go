package cpu

import (
	"6502profiler/memory"
	"fmt"
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

func TestADCZeroPageBCD(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x19
		c.Flags &= (^Flag_C)
		c.Flags |= Flag_D
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

	// adc $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x65, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC Zero page BCD",
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
		testProg:        []byte{0x75, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ADC Zero page X BCD",
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

// -------- EOR --------

func TestEORImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor #1
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x49, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 1)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x45, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 1)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x55, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x4d, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x5d, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
		c.Y = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor $1010,y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x59, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR absolute Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 1)
		c.Mem.Store(0x0012, 0x43)
		c.Mem.Store(0x0013, 0x15)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x41, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR index X indirect",
	}

	testSingleInstructionWithCase(t, c)
}

func TestEORIndirectY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 1)
		c.Mem.Store(0x0012, 0x40)
		c.Mem.Store(0x0013, 0x15)
		c.Y = 3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// eor ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x51, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR indirect Y",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- ORA --------

func TestORAImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora #4
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x09, 0x04, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 4)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x45, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ora Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 4)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x15, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 4)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x0d, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 4)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x1d, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 4)
		c.Y = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora $1010,y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x19, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA absolute Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 4)
		c.Mem.Store(0x0012, 0x43)
		c.Mem.Store(0x0013, 0x15)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// eor ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x01, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA index X indirect",
	}

	testSingleInstructionWithCase(t, c)
}

func TestORAIndirectY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 4)
		c.Mem.Store(0x0012, 0x40)
		c.Mem.Store(0x0013, 0x15)
		c.Y = 3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF7
	}

	// ora ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x11, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ORA indirect Y",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- AND --------

func TestANDImmediate(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and #1
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x29, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND immediate",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 1)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x25, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x0012, 1)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x35, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x2d, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x3d, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "EOR absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDAbsoluteY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1012, 1)
		c.Y = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and $1010,y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x39, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND absolute Y",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDIdxXIndirect(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 1)
		c.Mem.Store(0x0012, 0x43)
		c.Mem.Store(0x0013, 0x15)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// eor ($10, x)
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x21, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND index X indirect",
	}

	testSingleInstructionWithCase(t, c)
}

func TestANDIndirectY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
		c.Mem.Store(0x1543, 1)
		c.Mem.Store(0x0012, 0x40)
		c.Mem.Store(0x0013, 0x15)
		c.Y = 3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// and ($12), y
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x31, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "AND indirect Y",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- logical ops --------

func TestLsr(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x03
	res := Lsr(cpu, val)

	if res != 1 {
		t.Fatal("Lsr does not work")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Lsr does not set carry")
	}
}

func TestLsrNoCarry(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x02
	res := Lsr(cpu, val)

	if res != 1 {
		t.Fatal("Lsr does not work")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Lsr does set carry wrongly")
	}
}

func TestLsrZero(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x01
	res := Lsr(cpu, val)

	if res != 0 {
		t.Fatal("Lsr does not work")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Lsr does not set carry")
	}
}

func TestAsl(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x40
	res := Asl(cpu, val)

	if res != 0x80 {
		t.Fatal("Asl does not work")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Asl does set carry wrongly")
	}

	res = Asl(cpu, res)

	if res != 0x00 {
		t.Fatal("Asl does not work")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Asl does not set carry")
	}
}

func TestRol(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x80
	res := Rol(cpu, val)

	if res != 0x00 {
		t.Fatal("Rol does not work")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Rol does not set carry")
	}

	res = Rol(cpu, res)

	if res != 0x01 {
		t.Fatal("Rol does not work")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Rol does set carry wrongly")
	}
}

func TestRor(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	var val uint8 = 0x01
	res := Ror(cpu, val)

	if res != 0x00 {
		t.Fatal("Ror does not work")
	}

	if (cpu.Flags & Flag_C) == 0 {
		t.Fatal("Ror does not set carry")
	}

	res = Ror(cpu, res)

	if res != 0x80 {
		t.Fatal("Ror does not work")
	}

	if (cpu.Flags & Flag_C) != 0 {
		t.Fatal("Ror does set carry wrongly")
	}
}

func TestDec(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))

	res := Dec(cpu, 0x02)
	if res != 1 {
		t.Fatal("Dec does not work")
	}
}

func TestBit(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(1024))
	cpu.A = 0x0F

	cpu.bitBase(0xF0)
	if (cpu.Flags & Flag_N) == 0 {
		t.Fatal("Bit does not set negative flag")
	}

	if (cpu.Flags & Flag_V) == 0 {
		t.Fatal("Bit does not set overflow flag")
	}

	if (cpu.Flags & Flag_Z) == 0 {
		t.Fatal("Bit does not set zero flag")
	}
}

// -------- INC --------

func TestINCImplied(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF4
	}

	// inc
	// brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x1A, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INC A",
	}

	testSingleInstructionWithCase(t, c)
}

func TestINCZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0xF3)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0xF4
	}

	// inc $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xE6, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INC Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestINCZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0xF3)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0xF4
	}

	// inc $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xF6, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INC Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestINCAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0xF3)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0xF4
	}

	// inc $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xEE, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INC absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestINCAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0xF3)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0xF4
	}

	// inc $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xFE, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "INC absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- DEC --------

func TestDECImplied(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0xF3
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0xF2
	}

	// dec
	// brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x3A, 0x01, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEC A",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDECZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0xF3)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0xF2
	}

	// dec $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xC6, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEC Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDECZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0xF3)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0xF2
	}

	// dec $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xD6, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEC Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDECAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0xF3)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0xF2
	}

	// dec $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xCE, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEC absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestDECAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0xF3)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0xF2
	}

	// dec $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xDE, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "DEC absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- BIT --------

func TestBITZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x0F
		c.Mem.Store(0x0012, 0xF0)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_V) == 0 {
			return false
		}

		return true
	}

	// bit $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x24, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BIT Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestBITAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x0F
		c.Mem.Store(0x1012, 0xF0)
	}

	verifier := func(c *CPU6502) bool {
		if (c.Flags & Flag_N) == 0 {
			return false
		}

		if (c.Flags & Flag_Z) == 0 {
			return false
		}

		if (c.Flags & Flag_V) == 0 {
			return false
		}

		return true
	}

	// bit $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x2C, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "BIT absolute",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TAX --------

func TestTAX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x0F
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.X == 0x0F
	}

	// tax
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xAA, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TAX",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TAX --------

func TestTXA(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 0x0F
		c.A = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x0F
	}

	// txa
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x8A, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TXA",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TAY --------

func TestTAY(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x0F
		c.Y = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Y == 0x0F
	}

	// tay
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xA8, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TAY",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TYA --------

func TestTYA(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Y = 0x0F
		c.A = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x0F
	}

	// tya
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x98, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TYA",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TXS --------

func TestTXS(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.X = 0x0F
		c.SP = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.SP == 0x0F
	}

	// txs
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x9A, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TXS",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- TSX --------

func TestTSX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.SP = 0x0F
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.X == 0x0F
	}

	// tsx
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0xBA, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "TSX",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- ASL --------

func TestASLImplied(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x02
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x04
	}

	// asl
	// brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x0A, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ASL A",
	}

	testSingleInstructionWithCase(t, c)
}

func TestASLZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x04
	}

	// asl $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x06, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ASL Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestASLZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x04
	}

	// asl $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x16, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ASL Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestASLAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x04
	}

	// asl $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x0E, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ASL absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestASLAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x04
	}

	// asl $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x1E, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ASL absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- LSR --------

func TestLSRImplied(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x02
	}

	verifier := func(c *CPU6502) bool {
		return c.A == 0x01
	}

	// lsr
	// brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x4A, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LSR A",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLSRZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x01
	}

	// lsr $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x46, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LSR Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLSRZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x01
	}

	// lsr $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x56, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LSR Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLSRAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x01
	}

	// lsr $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x5E, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LSR absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestLSRAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x01
	}

	// lsr $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x5E, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "LSR absolute X",
	}

	testSingleInstructionWithCase(t, c)
}

// -------- ROL --------

func TestROLImplied(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.A = 0x02
		c.Flags |= Flag_C
	}

	verifier := func(c *CPU6502) bool {
		fmt.Println(c.A)
		return c.A == 0x05
	}

	// rol
	// brk
	c := InstructionTestCase{
		model:           Model65C02,
		testProg:        []byte{0x2A, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ROL A",
	}

	testSingleInstructionWithCase(t, c)
}

func TestROLZeroPage(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
		c.Flags |= Flag_C
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x05
	}

	// rol $12
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x26, 0x12, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ROL Zero page",
	}

	testSingleInstructionWithCase(t, c)
}

func TestROLZeroPageX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x0012, 0x02)
		c.Flags |= Flag_C
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x0012) == 0x05
	}

	// rol $10,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x36, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ROL Zero page X",
	}

	testSingleInstructionWithCase(t, c)
}

func TestROLAbsolute(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
		c.Flags |= Flag_C
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x05
	}

	// rol $1012
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x2E, 0x12, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ROL absolute",
	}

	testSingleInstructionWithCase(t, c)
}

func TestROLAbsoluteX(t *testing.T) {
	arranger := func(c *CPU6502) {
		c.Mem.Store(0x1012, 0x02)
		c.Flags |= Flag_C
		c.X = 2
	}

	verifier := func(c *CPU6502) bool {
		return c.Mem.Load(0x1012) == 0x05
	}

	// rol $1010,x
	// brk
	c := InstructionTestCase{
		model:           Model6502,
		testProg:        []byte{0x3E, 0x10, 0x10, 0x00},
		arranger:        arranger,
		verifier:        verifier,
		instructionName: "ROL absolute X",
	}

	testSingleInstructionWithCase(t, c)
}
