package main

const Flag_N uint8 = 0x80
const Flag_V uint8 = 0x40
const Flag_B uint8 = 0x10
const Flag_D uint8 = 0x08
const Flag_I uint8 = 0x04
const Flag_Z uint8 = 0x02
const Flag_C uint8 = 0x0

type CpuModel uint8

const Model6502 CpuModel = 0x00
const Model65C02 CpuModel = 0x01

type Memory interface {
	get(address uint16) uint8
	set(address uint16, b uint8)
}

type CPU6502 struct {
	PC         uint16
	SP         uint8
	A          uint8
	X          uint8
	Y          uint8
	Flags      uint8
	model      CpuModel
	cycleCount uint64
	mem        Memory
}

func New6502(m CpuModel) *CPU6502 {
	res := &CPU6502{
		PC:         0x0000,
		SP:         0xFF,
		A:          0,
		X:          0,
		Y:          0,
		Flags:      0x00,
		model:      m,
		cycleCount: 0,
	}

	return res
}

func (c *CPU6502) Init(m Memory) {
	c.SP = 0xFF
	c.cycleCount = 0
	c.mem = m
}

func (c *CPU6502) Run(startAddress uint16) {
	c.PC = startAddress
	var err error

	for err = nil; err == nil; {
		err = c.executeInstruction()
	}
}

func (c *CPU6502) executeInstruction() error {
	return nil
}

func main() {
	cpu := New6502(Model6502)
	cpu.Init(nil)
	cpu.Run(0x1000)
}
