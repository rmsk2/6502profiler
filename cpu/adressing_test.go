package cpu

import (
	"6502profiler/memory"
	"testing"
)

func TestRelativeOverflow(t *testing.T) {
	calcOffset := func(a uint16, o uint8) uint16 {
		offset2 := int16(int8(o))
		return uint16(int16(a) + offset2)
	}

	if calcOffset(32769, 0xFF) != 32768 {
		t.Fatal("Offset calculation does not work")
	}

	if calcOffset(10, 0xFF) != 9 {
		t.Fatal("Offset calculation does not work")
	}

	if calcOffset(0, 0xFF) != 65535 {
		t.Fatal("Offset calculation does not work")
	}

	if calcOffset(0xFFFF, 1) != 0 {
		t.Fatal("Offset calculation does not work")
	}
}

func TestGetAddrAbsolute(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x42, 0x24}, 0x0000)
	cpu.PC = 0x0000

	if cpu.getAddrAbsolute() != 0x2442 {
		t.Fatal("Absolute addressing does not work")
	}
}

func TestGetAddrZeroPage(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x42}, 0x0000)
	cpu.PC = 0x0000

	if cpu.getAddrZeroPage() != 0x42 {
		t.Fatal("Zero page addressing does not work")
	}
}

func TestGetAddrAbsoluteY(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x42, 0x24}, 0x0000)
	cpu.PC = 0x0000
	cpu.Y = 0x08

	res, _ := cpu.getAddrAbsoluteY()

	if res != 0x244A {
		t.Fatal("Absolute Y addressing does not work")
	}
}

func TestGetAddrAbsoluteX(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x52, 0x25}, 0x0000)
	cpu.PC = 0x0000
	cpu.X = 0x08

	res, _ := cpu.getAddrAbsoluteX()

	if res != 0x255A {
		t.Fatal("Absolute X addressing does not work")
	}
}

func TestGetAddrZeroPageY(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x17}, 0x0000)
	cpu.PC = 0x0000
	cpu.Y = 0xFF

	if cpu.getAddrZeroPageY() != 0x0016 {
		t.Fatal("Zero page Y addressing does not work")
	}
}

func TestGetAddrZeroPageX(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xFF}, 0x0000)
	cpu.PC = 0x0000
	cpu.X = 0xFF

	if cpu.getAddrZeroPageX() != 0x00FE {
		t.Fatal("Zero page X addressing does not work")
	}
}

func TestGetAddrIndirect(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x02, 0x00, 0x89, 0xAF}, 0x0000)
	cpu.PC = 0x0000

	if cpu.getAddrIndirect() != 0xAF89 {
		t.Fatal("Indirect addressing does not work")
	}
}

func TestGetAddrJmIndirect65C02(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xFF, 0x30}, 0x0000)
	cpu.Mem.Store(0x30FF, 0x80)
	cpu.Mem.Store(0x3000, 0x40)
	cpu.Mem.Store(0x3100, 0x50)

	addr := cpu.getAddrIndirect()

	if addr != 0x5080 {
		t.Fatal("Indirect JMP does not work")
	}
}

func TestGetAddrRelativeNegative(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xa9, 0x00, 0xf0, 0xfe}, 0x0000)
	cpu.PC = 0x0003

	res, _ := cpu.getAddrRelative()

	if res != 0x0002 {
		t.Fatal("Relative addressing for a negative offset does not work")
	}
}

func TestGetAddrRelativePositive(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xa9, 0x00, 0xf0, 0x00, 0xea}, 0x0000)
	cpu.PC = 0x0003

	res, _ := cpu.getAddrRelative()

	if res != 0x0004 {
		t.Fatal("Relative addressing for a positive	offset does not work")
	}
}

func TestGetAddrIndirectIdxY(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x12}, 0x0000)
	cpu.Mem.Store(0x0012, 0x78)
	cpu.Mem.Store(0x0013, 0x56)
	cpu.Y = 2
	cpu.PC = 0x0000

	res, _ := cpu.getAddrIndirectIdxY()

	if res != 0x567A {
		t.Fatal("Indirect addressing with index Y does not work")
	}
}

func TestGetAddrIdxIndirectX(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x12}, 0x0000)
	cpu.Mem.Store(0x0012, 0x78)
	cpu.Mem.Store(0x0013, 0x56)
	cpu.Mem.Store(0x0014, 0xBC)
	cpu.Mem.Store(0x0015, 0x9A)
	cpu.X = 2
	cpu.PC = 0x0000

	res := cpu.getAddrIdxIndirectX()

	if res != 0x9ABC {
		t.Fatal("indexed with X indirect addressing does not work")
	}
}

// This test verifies compatibility to a bug present in the original 6502
func TestGetAddrJmIndirectJmp6502(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xFF, 0x30}, 0x0000)
	cpu.Mem.Store(0x30FF, 0x80)
	cpu.Mem.Store(0x3000, 0x40)
	cpu.Mem.Store(0x3100, 0x50)

	addr := cpu.getAddrIndirectJmp6502()

	if addr != 0x4080 {
		t.Fatal("indirect JMP works too well for an original 6502")
	}
}

func TestGetAddrJmIndirectJmp6502Correct(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0xFE, 0x30}, 0x0000)
	cpu.Mem.Store(0x30FE, 0x80)
	cpu.Mem.Store(0x30FF, 0x50)

	addr := cpu.getAddrIndirectJmp6502()

	if addr != 0x5080 {
		t.Fatal("indirect JMP does not work")
	}
}

func TestGetAddrIdxIndirect65C02(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x57, 0x13}, 0x0000)
	cpu.Mem.Store(0x1456, 0xCD)
	cpu.Mem.Store(0x1457, 0xAB)
	cpu.X = 0xFF
	cpu.PC = 0x0000

	res := cpu.getAddrIdxIndirect65C02()

	if res != 0xABCD {
		t.Fatal("Indexed indirect addressing (65C02) does not work")
	}
}

func TestGetAddrZp65C02(t *testing.T) {
	cpu := New6502(Model6502)
	cpu.Init(memory.NewLinearMemory(16384))
	cpu.CopyToMem([]byte{0x12}, 0x0000)
	cpu.Mem.Store(0x0012, 0x78)
	cpu.Mem.Store(0x0013, 0x56)
	cpu.PC = 0x0000

	res := cpu.getAddrZp65C02()

	if res != 0x5678 {
		t.Fatal("65C02 Zero-page addressing does not work")
	}
}
