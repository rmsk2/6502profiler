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

type CPU6502 struct {
	PC         uint16
	SP         uint8
	A          uint8
	X          uint8
	Y          uint8
	Flags      uint8
	model      CpuModel
	cycleCount uint64
	mem        memory.Memory
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

//type InstructionFunc func(c *CPU6502, operAddr uint16, operand uint8) (uint64, bool)
//type FetchOperFunc func(c *CPU6502) uint8

func (c *CPU6502) Init(m memory.Memory) {
	c.mem = m
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
		c.mem.Store(loadAddress, j)
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

func (c *CPU6502) executeInstruction() (uint64, bool) {
	//panic(fmt.Sprintf("Shit at %d", c.PC))
	return 0, true
}
