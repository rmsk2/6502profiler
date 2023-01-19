package cpu

import (
	"6502profiler/memory"
	"fmt"
	"os"
)

const Flag_N uint8 = 0x80
const Flag_V uint8 = 0x40
const Flag_B uint8 = 0x10
const Flag_D uint8 = 0x08
const Flag_I uint8 = 0x04
const Flag_Z uint8 = 0x02
const Flag_C uint8 = 0x01

type CpuModel uint8

const Model6502 CpuModel = 0x00
const Model65C02 CpuModel = 0x01

type execFunc func(c *CPU6502) (uint64, bool)

type CPU6502 struct {
	PC         uint16
	SP         uint8
	A          uint8
	X          uint8
	Y          uint8
	Flags      uint8
	model      CpuModel
	cycleCount uint64
	Mem        memory.Memory
	opCodes    map[byte]execFunc
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
		opCodes:    make(map[uint8]execFunc),
	}

	res.opCodes[0xA9] = (*CPU6502).ldaImmediate
	res.opCodes[0xAD] = (*CPU6502).ldaAbsolute
	res.opCodes[0x8D] = (*CPU6502).staAbsolute
	res.opCodes[0x00] = (*CPU6502).brk

	return res
}

func (c *CPU6502) NumCycles() uint64 {
	return c.cycleCount
}

func (c *CPU6502) Init(m memory.Memory) {
	c.Mem = m
	c.SP = 0xFF
}

func (c *CPU6502) LoadAndRun(fileName string) (err error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("unable to load binary: %v", err)
	}

	if len(data) < 3 {
		return fmt.Errorf("no program data found")
	}

	var loadAddress uint16 = uint16(data[1])*256 + uint16(data[0])
	startAddress := loadAddress

	// Recover from panic created by memory access
	defer func() {
		if res := recover(); res != nil {
			// Use named return value to return a value after handling the panic
			err = fmt.Errorf("error loading 6502 program: %v", res)
		}
	}()

	for _, j := range data[2:] {
		c.Mem.Store(loadAddress, j)
		loadAddress++ // This can overflow
	}

	err = c.Run(startAddress)

	return err
}

func (c *CPU6502) Run(startAddress uint16) (err error) {
	var cyclesUsed uint64
	err = nil

	// Recover from panic created by an instruction
	defer func() {
		if res := recover(); res != nil {
			// Use named return value to return a value after handling the panic
			err = fmt.Errorf("error running 6502 program: %v", res)
		}
	}()

	c.PC = startAddress
	c.cycleCount = 0

	for halt := false; !halt; {
		cyclesUsed, halt = c.executeInstruction()
		if !halt {
			c.cycleCount += cyclesUsed
		}
	}

	return err
}

func (c *CPU6502) getAddrAbsolute() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return addr
}

func (c *CPU6502) ldaBase(value uint8) bool {
	c.A = value

	if c.A == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &^= Flag_Z
	}

	if (c.A & 0x80) != 0 {
		c.Flags |= Flag_N
	} else {
		c.Flags &^= Flag_N
	}

	return false
}

func (c *CPU6502) ldaImmediate() (uint64, bool) {
	stop := c.ldaBase(c.Mem.Load(c.PC))
	c.PC++

	return 2, stop
}

func (c *CPU6502) ldaAbsolute() (uint64, bool) {
	stop := c.ldaBase(c.Mem.Load(c.getAddrAbsolute()))
	c.PC++

	return 4, stop
}

func (c *CPU6502) staAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.A)
	c.PC++

	return 4, false
}

func (c *CPU6502) brk() (uint64, bool) {
	return 7, true
}

func (c *CPU6502) executeInstruction() (uint64, bool) {
	opCode := c.Mem.Load(c.PC)
	instruction, ok := c.opCodes[opCode]
	if !ok {
		panic(fmt.Sprintf("Illegal opcode $%x at $%x", opCode, c.PC))
	}

	c.PC++

	return instruction(c)
}
