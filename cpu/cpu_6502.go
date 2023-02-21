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

	// BPL
	res.opCodes[0x10] = (*CPU6502).bpl
	// BMI
	res.opCodes[0x30] = (*CPU6502).bmi
	// BEQ
	res.opCodes[0xF0] = (*CPU6502).beq
	// BNE
	res.opCodes[0xD0] = (*CPU6502).bne
	// BCC
	res.opCodes[0x90] = (*CPU6502).bcc
	// BCS
	res.opCodes[0xB0] = (*CPU6502).bcs
	// BVC
	res.opCodes[0x50] = (*CPU6502).bvc
	// BVS
	res.opCodes[0x70] = (*CPU6502).bvs

	// CPY
	res.opCodes[0xC0] = (*CPU6502).cpyImmediate
	res.opCodes[0xC4] = (*CPU6502).cpyZeroPage
	res.opCodes[0xCC] = (*CPU6502).cpyAbsolute

	// CPX
	res.opCodes[0xE0] = (*CPU6502).cpxImmediate
	res.opCodes[0xE4] = (*CPU6502).cpxZeroPage
	res.opCodes[0xEC] = (*CPU6502).cpxAbsolute

	// DEY
	res.opCodes[0x88] = (*CPU6502).dey
	// INY
	res.opCodes[0xC8] = (*CPU6502).iny

	// DEX
	res.opCodes[0xCA] = (*CPU6502).dex
	// INX
	res.opCodes[0xE8] = (*CPU6502).inx

	// LDA
	res.opCodes[0xA9] = (*CPU6502).ldaImmediate
	res.opCodes[0xAD] = (*CPU6502).ldaAbsolute
	res.opCodes[0xB9] = (*CPU6502).ldaAbsoluteY
	res.opCodes[0xBD] = (*CPU6502).ldaAbsoluteX
	res.opCodes[0xB1] = (*CPU6502).ldaIndIdxY
	res.opCodes[0xA5] = (*CPU6502).ldaZeroPage
	res.opCodes[0xB5] = (*CPU6502).ldaZeroPageIdxX
	res.opCodes[0xA1] = (*CPU6502).ldaIdxIndirectX

	// LDX
	res.opCodes[0xA2] = (*CPU6502).ldxImmediate
	res.opCodes[0xBE] = (*CPU6502).ldxAbsoluteY
	res.opCodes[0xAE] = (*CPU6502).ldxAbsolute
	res.opCodes[0xA6] = (*CPU6502).ldxZeroPage
	res.opCodes[0xB6] = (*CPU6502).ldxZeroPageIdxY

	// LDY
	res.opCodes[0xAC] = (*CPU6502).ldyAbsolute
	res.opCodes[0xA0] = (*CPU6502).ldyImmediate
	res.opCodes[0xBC] = (*CPU6502).ldyAbsoluteX
	res.opCodes[0xA4] = (*CPU6502).ldyZeroPage
	res.opCodes[0xB4] = (*CPU6502).ldyZeroPageIdxX

	// STA
	res.opCodes[0x8D] = (*CPU6502).staAbsolute
	res.opCodes[0x99] = (*CPU6502).staAbsoluteY
	res.opCodes[0x85] = (*CPU6502).staZeroPage
	res.opCodes[0x9D] = (*CPU6502).staAbsoluteX
	res.opCodes[0x95] = (*CPU6502).staZeroPageX
	res.opCodes[0x91] = (*CPU6502).staIndirectY
	res.opCodes[0x81] = (*CPU6502).staXIndirect

	// STX
	res.opCodes[0x86] = (*CPU6502).stxZeroPage
	res.opCodes[0x96] = (*CPU6502).stxZeroPageY
	res.opCodes[0x8E] = (*CPU6502).stxAbsolute

	// STY
	res.opCodes[0x84] = (*CPU6502).styZeroPage
	res.opCodes[0x94] = (*CPU6502).styZeroPageX
	res.opCodes[0x8C] = (*CPU6502).styAbsolute

	// CMP
	res.opCodes[0xC9] = (*CPU6502).cmpImmediate
	res.opCodes[0xC5] = (*CPU6502).cmpZeroPage
	res.opCodes[0xD5] = (*CPU6502).cmpZeroPageX
	res.opCodes[0xCD] = (*CPU6502).cmpAbsolute
	res.opCodes[0xDD] = (*CPU6502).cmpAbsoluteX
	res.opCodes[0xD9] = (*CPU6502).cmpAbsoluteY
	res.opCodes[0xC1] = (*CPU6502).cmpIdxXIndirect
	res.opCodes[0xD1] = (*CPU6502).cmpIndIdxY

	// ADC
	res.opCodes[0x69] = (*CPU6502).addImmediate
	res.opCodes[0x65] = (*CPU6502).addZeroPage
	res.opCodes[0x75] = (*CPU6502).addZeroPageX
	res.opCodes[0x6D] = (*CPU6502).addAbsolute
	res.opCodes[0x7D] = (*CPU6502).addAbsoluteX
	res.opCodes[0x79] = (*CPU6502).addAbsoluteY
	res.opCodes[0x71] = (*CPU6502).addIndirectIdxY
	res.opCodes[0x61] = (*CPU6502).addIdxXIndirect

	// SBC
	res.opCodes[0xE9] = (*CPU6502).subImmediate
	res.opCodes[0xE5] = (*CPU6502).subZeroPage
	res.opCodes[0xF5] = (*CPU6502).subZeroPageX
	res.opCodes[0xED] = (*CPU6502).subAbsolute
	res.opCodes[0xFD] = (*CPU6502).subAbsoluteX
	res.opCodes[0xF9] = (*CPU6502).subAbsoluteY
	res.opCodes[0xF1] = (*CPU6502).subIndirectIdxY
	res.opCodes[0xE1] = (*CPU6502).subIdxXIndirect

	// EOR
	res.opCodes[0x49] = (*CPU6502).eorImmediate
	res.opCodes[0x45] = (*CPU6502).eorZeroPage
	res.opCodes[0x55] = (*CPU6502).eorZeroPageX
	res.opCodes[0x4D] = (*CPU6502).eorAbsolute
	res.opCodes[0x5D] = (*CPU6502).eorAbsoluteX
	res.opCodes[0x59] = (*CPU6502).eorAbsoluteY
	res.opCodes[0x41] = (*CPU6502).eorIdxIndirect
	res.opCodes[0x51] = (*CPU6502).eorIndirectIdxY

	// ORA
	res.opCodes[0x09] = (*CPU6502).oraImmediate
	res.opCodes[0x05] = (*CPU6502).oraZeroPage
	res.opCodes[0x15] = (*CPU6502).oraZeroPageX
	res.opCodes[0x0D] = (*CPU6502).oraAbsolute
	res.opCodes[0x1D] = (*CPU6502).oraAbsoluteX
	res.opCodes[0x19] = (*CPU6502).oraAbsoluteY
	res.opCodes[0x01] = (*CPU6502).oraIdxIndirect
	res.opCodes[0x11] = (*CPU6502).oraIndirectIdxY

	// AND
	res.opCodes[0x29] = (*CPU6502).andImmediate
	res.opCodes[0x25] = (*CPU6502).andZeroPage
	res.opCodes[0x35] = (*CPU6502).andZeroPageX
	res.opCodes[0x2D] = (*CPU6502).andAbsolute
	res.opCodes[0x3D] = (*CPU6502).andAbsoluteX
	res.opCodes[0x39] = (*CPU6502).andAbsoluteY
	res.opCodes[0x21] = (*CPU6502).andIdxIndirect
	res.opCodes[0x31] = (*CPU6502).andIndirectIdxY

	// INC
	res.opCodes[0xE6] = (*CPU6502).incZeroPage
	res.opCodes[0xF6] = (*CPU6502).incZeroPageX
	res.opCodes[0xEE] = (*CPU6502).incAbsolute
	res.opCodes[0xFE] = (*CPU6502).incAbsoluteX

	// DEC
	res.opCodes[0xC6] = (*CPU6502).decZeroPage
	res.opCodes[0xD6] = (*CPU6502).decZeroPageX
	res.opCodes[0xCE] = (*CPU6502).decAbsolute
	res.opCodes[0xDE] = (*CPU6502).decAbsoluteX

	// ASL
	res.opCodes[0x0a] = (*CPU6502).asl
	res.opCodes[0x06] = (*CPU6502).aslZeroPage
	res.opCodes[0x16] = (*CPU6502).aslZeroPageX
	res.opCodes[0x0E] = (*CPU6502).aslAbsolute
	res.opCodes[0x1E] = (*CPU6502).aslAbsoluteX

	// LSR
	res.opCodes[0x4A] = (*CPU6502).lsr
	res.opCodes[0x46] = (*CPU6502).lsrZeroPage
	res.opCodes[0x56] = (*CPU6502).lsrZeroPageX
	res.opCodes[0x4E] = (*CPU6502).lsrAbsolute
	res.opCodes[0x5E] = (*CPU6502).lsrAbsoluteX

	// ROL
	res.opCodes[0x2A] = (*CPU6502).rol
	res.opCodes[0x26] = (*CPU6502).rolZeroPage
	res.opCodes[0x36] = (*CPU6502).rolZeroPageX
	res.opCodes[0x2E] = (*CPU6502).rolAbsolute
	res.opCodes[0x3E] = (*CPU6502).rolAbsoluteX

	// ROR
	res.opCodes[0x6A] = (*CPU6502).ror
	res.opCodes[0x66] = (*CPU6502).rorZeroPage
	res.opCodes[0x76] = (*CPU6502).rorZeroPageX
	res.opCodes[0x6E] = (*CPU6502).rorAbsolute
	res.opCodes[0x7E] = (*CPU6502).rorAbsoluteX

	// BIT
	res.opCodes[0x24] = (*CPU6502).bitZeroPage
	res.opCodes[0x2C] = (*CPU6502).bitAbsolute

	// JSR
	res.opCodes[0x20] = (*CPU6502).jsr

	// RTS
	res.opCodes[0x60] = (*CPU6502).rts

	// JMP
	res.opCodes[0x4C] = (*CPU6502).jmp
	if m == Model6502 {
		res.opCodes[0x6c] = (*CPU6502).jmpIndirect6502
	}

	// PHA
	res.opCodes[0x48] = (*CPU6502).pha
	// PLA
	res.opCodes[0x68] = (*CPU6502).pla
	// PLP
	res.opCodes[0x28] = (*CPU6502).plp
	// PHP
	res.opCodes[0x08] = (*CPU6502).php

	// TAX
	res.opCodes[0xAA] = (*CPU6502).tax

	// TXA
	res.opCodes[0x8A] = (*CPU6502).txa

	// TAY
	res.opCodes[0xA8] = (*CPU6502).tay

	// TYA
	res.opCodes[0x98] = (*CPU6502).tya

	// TXS
	res.opCodes[0x9A] = (*CPU6502).txs

	// TSX
	res.opCodes[0xBA] = (*CPU6502).tsx

	// Flag stuff
	res.opCodes[0x18] = (*CPU6502).clc
	res.opCodes[0xD8] = (*CPU6502).cld
	res.opCodes[0x58] = (*CPU6502).cli
	res.opCodes[0xB8] = (*CPU6502).clv
	res.opCodes[0x38] = (*CPU6502).sec
	res.opCodes[0xF8] = (*CPU6502).sed
	res.opCodes[0x78] = (*CPU6502).sei

	// BRK
	res.opCodes[0x00] = func(c *CPU6502) (uint64, bool) {
		return 7, true
	}

	// NOP
	res.opCodes[0xEA] = func(c *CPU6502) (uint64, bool) {
		return 2, false
	}

	if m == Model65C02 {
		// New instructions
		res.opCodes[0x80] = (*CPU6502).bra
		res.opCodes[0x64] = (*CPU6502).stzZeroPage
		res.opCodes[0x74] = (*CPU6502).stzZeroPageX
		res.opCodes[0x9C] = (*CPU6502).stzAbsolute
		res.opCodes[0x9E] = (*CPU6502).stzAbsoluteX
		res.opCodes[0xDA] = (*CPU6502).phx
		res.opCodes[0xFA] = (*CPU6502).plx
		res.opCodes[0x5A] = (*CPU6502).phy
		res.opCodes[0x7A] = (*CPU6502).ply
		res.opCodes[0x14] = (*CPU6502).trbZeroPage
		res.opCodes[0x1c] = (*CPU6502).trbAbsolute
		res.opCodes[0x04] = (*CPU6502).tsbZeroPage
		res.opCodes[0x0c] = (*CPU6502).tsbAbsolute

		res.opCodes[0x0f] = (*CPU6502).bbr0
		res.opCodes[0x1f] = (*CPU6502).bbr1
		res.opCodes[0x2f] = (*CPU6502).bbr2
		res.opCodes[0x3f] = (*CPU6502).bbr3
		res.opCodes[0x4f] = (*CPU6502).bbr4
		res.opCodes[0x5f] = (*CPU6502).bbr5
		res.opCodes[0x6f] = (*CPU6502).bbr6
		res.opCodes[0x7f] = (*CPU6502).bbr7

		res.opCodes[0x8f] = (*CPU6502).bbs0
		res.opCodes[0x9f] = (*CPU6502).bbs1
		res.opCodes[0xaf] = (*CPU6502).bbs2
		res.opCodes[0xbf] = (*CPU6502).bbs3
		res.opCodes[0xcf] = (*CPU6502).bbs4
		res.opCodes[0xdf] = (*CPU6502).bbs5
		res.opCodes[0xef] = (*CPU6502).bbs6
		res.opCodes[0xff] = (*CPU6502).bbs7

		res.opCodes[0x07] = (*CPU6502).rmb0
		res.opCodes[0x17] = (*CPU6502).rmb1
		res.opCodes[0x27] = (*CPU6502).rmb2
		res.opCodes[0x37] = (*CPU6502).rmb3
		res.opCodes[0x47] = (*CPU6502).rmb4
		res.opCodes[0x57] = (*CPU6502).rmb5
		res.opCodes[0x67] = (*CPU6502).rmb6
		res.opCodes[0x77] = (*CPU6502).rmb7

		res.opCodes[0x87] = (*CPU6502).smb0
		res.opCodes[0x97] = (*CPU6502).smb1
		res.opCodes[0xa7] = (*CPU6502).smb2
		res.opCodes[0xb7] = (*CPU6502).smb3
		res.opCodes[0xc7] = (*CPU6502).smb4
		res.opCodes[0xd7] = (*CPU6502).smb5
		res.opCodes[0xe7] = (*CPU6502).smb6
		res.opCodes[0xf7] = (*CPU6502).smb7

		// New addressing modes for exisiting instructions
		res.opCodes[0x6c] = (*CPU6502).jmpIndirect65C02
		res.opCodes[0x7c] = (*CPU6502).jmpIndexXIndirect
		res.opCodes[0x1a] = (*CPU6502).inc65C02
		res.opCodes[0x3a] = (*CPU6502).dec65C02

		res.opCodes[0x89] = (*CPU6502).bitImmediate
		res.opCodes[0x34] = (*CPU6502).bitZeroPageX
		res.opCodes[0x3C] = (*CPU6502).bitAbsoluteX

		res.opCodes[0x72] = (*CPU6502).addIndirect
		res.opCodes[0xF2] = (*CPU6502).subIndirect
		res.opCodes[0x32] = (*CPU6502).andIndirect
		res.opCodes[0x52] = (*CPU6502).eorIndirect
		res.opCodes[0x12] = (*CPU6502).oraIndirect
		res.opCodes[0xd2] = (*CPU6502).cmpIndirect
		res.opCodes[0xb2] = (*CPU6502).ldaIndirect
		res.opCodes[0x92] = (*CPU6502).staIndirect

		// Different cycle count when compared with a 6502. ADD and SBC are missing here
		// because the differences in cycle count for these instructions are implemented in
		// the ADD and SBC routines directly.
		res.opCodes[0x1e] = (*CPU6502).aslAbsoluteX65C02
		res.opCodes[0x5e] = (*CPU6502).lsrAbsoluteX65C02
		res.opCodes[0x3e] = (*CPU6502).rolAbsoluteX65C02
		res.opCodes[0x7e] = (*CPU6502).rorAbsoluteX65C02
	}

	return res
}

func (c *CPU6502) Reset() {
	c.cycleCount = 0
	c.Flags = 0
	c.X = 0
	c.A = 0
	c.Y = 0
	c.PC = 0
	c.SP = 0xFF
	c.Mem.ClearStatistics()
}

func (c *CPU6502) NumCycles() uint64 {
	return c.cycleCount
}

func (c *CPU6502) Init(m memory.Memory) {
	c.Mem = m
	c.SP = 0xFF
}

func (c *CPU6502) LoadAndRun(fileName string) (loadAddress uint16, progLen uint16, err error) {
	loadAddress, progLen, err = c.Load(fileName)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to run program: %v", err)
	}

	return loadAddress, progLen, c.Run(loadAddress)
}

func (c *CPU6502) Load(fileName string) (loadAddress uint16, progLen uint16, err error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return 0, 0, fmt.Errorf("unable to load binary: %v", err)
	}

	if len(data) < 3 {
		return 0, 0, fmt.Errorf("no program data found")
	}

	loadAddress = uint16(data[1])*256 + uint16(data[0])
	c.CopyToMem(data[2:], loadAddress)

	return loadAddress, uint16(len(data) - 2), nil
}

func (c *CPU6502) CopyToMem(binary []byte, startAddress uint16) (err error) {
	// Recover from panic created by memory access
	defer func() {
		if res := recover(); res != nil {
			// Use named return value to return a value after handling the panic
			err = fmt.Errorf("error copying to memory: %v", res)
		}
	}()

	copyAddress := startAddress

	for _, j := range binary {
		c.Mem.Store(copyAddress, j)
		copyAddress++ // This can overflow
	}

	return err
}

func (c *CPU6502) CopyFromMem(startAddress uint16, length uint16) (data []byte, err error) {
	// Recover from panic created by memory access
	defer func() {
		if res := recover(); res != nil {
			// Use named return value to return a value after handling the panic
			err = fmt.Errorf("error copying from memory: %v", res)
		}
	}()

	data = []byte{}
	copyAddress := startAddress
	var count uint16

	for count = 0; count < length; count++ {
		data = append(data, c.Mem.Load(copyAddress))
		copyAddress++
	}

	return data, err
}

func (c *CPU6502) CopyAndRun(program []byte, startAddress uint16) (err error) {
	err = c.CopyToMem(program, startAddress)
	if err != nil {
		return err
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

// -------- Helpers --------

func (c *CPU6502) nzFlags(v uint8) {
	if v == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &^= Flag_Z
	}

	if (v & 0x80) != 0 {
		c.Flags |= Flag_N
	} else {
		c.Flags &^= Flag_N
	}
}

// Stack functions. Stack is always in the area 0x100 - 0x1FF. The first address is
// 0x1FF. The stack grows downwards.
func (c *CPU6502) push(val uint8) {
	c.Mem.Store(0x100+uint16(c.SP), val)
	c.SP--
}

func (c *CPU6502) pop() uint8 {
	c.SP++
	return c.Mem.Load(0x100 + uint16(c.SP))
}

// ---------------------

func (c *CPU6502) executeInstruction() (uint64, bool) {
	opCode := c.Mem.Load(c.PC)
	instruction, ok := c.opCodes[opCode]
	if !ok {
		panic(fmt.Sprintf("Illegal opcode $%x at $%x", opCode, c.PC))
	}

	c.PC++

	return instruction(c)
}
