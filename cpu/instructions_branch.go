package cpu

func (c *CPU6502) branchOnFlagClear(flag uint8) (uint64, bool) {
	if (c.Flags & flag) != 0 {
		c.PC++ // Skip when flag is set
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

func (c *CPU6502) branchOnFlagSet(flag uint8) (uint64, bool) {
	if (c.Flags & flag) == 0 {
		c.PC++ // Skip when flag is clear
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

func (c *CPU6502) branchOnBitClear(bit uint8) (uint64, bool) {
	zpAddr, branchAddr, additionalCycle := c.getAddressesBitBranchRelative()
	if (c.Mem.Load(zpAddr) & bit) != 0 {
		c.PC++
		return 5, false
	}

	c.PC = branchAddr
	return 6 + additionalCycle, false
}

func (c *CPU6502) branchOnBitSet(bit uint8) (uint64, bool) {
	zpAddr, branchAddr, additionalCycle := c.getAddressesBitBranchRelative()
	if (c.Mem.Load(zpAddr) & bit) == 0 {
		c.PC++
		return 5, false
	}

	c.PC = branchAddr
	return 6 + additionalCycle, false
}

// -------- BPL --------

func (c *CPU6502) bpl() (uint64, bool) {
	return c.branchOnFlagClear(Flag_N)
}

// -------- BMI --------

func (c *CPU6502) bmi() (uint64, bool) {
	return c.branchOnFlagSet(Flag_N)
}

// -------- BNE --------

func (c *CPU6502) bne() (uint64, bool) {
	return c.branchOnFlagClear(Flag_Z)
}

// -------- BEQ --------

func (c *CPU6502) beq() (uint64, bool) {
	return c.branchOnFlagSet(Flag_Z)
}

// -------- BCC --------

func (c *CPU6502) bcc() (uint64, bool) {
	return c.branchOnFlagClear(Flag_C)
}

// -------- BCS--------

func (c *CPU6502) bcs() (uint64, bool) {
	return c.branchOnFlagSet(Flag_C)
}

// -------- BVC --------

func (c *CPU6502) bvc() (uint64, bool) {
	return c.branchOnFlagClear(Flag_V)
}

// -------- BVS--------

func (c *CPU6502) bvs() (uint64, bool) {
	return c.branchOnFlagSet(Flag_V)
}

// -------- BRA--------

func (c *CPU6502) bra() (uint64, bool) {
	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BBR0 - BBR7 --------

func (c *CPU6502) bbr0() (uint64, bool) {
	return c.branchOnBitClear(0x01)
}

func (c *CPU6502) bbr1() (uint64, bool) {
	return c.branchOnBitClear(0x02)
}

func (c *CPU6502) bbr2() (uint64, bool) {
	return c.branchOnBitClear(0x04)
}

func (c *CPU6502) bbr3() (uint64, bool) {
	return c.branchOnBitClear(0x08)
}

func (c *CPU6502) bbr4() (uint64, bool) {
	return c.branchOnBitClear(0x10)
}

func (c *CPU6502) bbr5() (uint64, bool) {
	return c.branchOnBitClear(0x20)
}

func (c *CPU6502) bbr6() (uint64, bool) {
	return c.branchOnBitClear(0x40)
}

func (c *CPU6502) bbr7() (uint64, bool) {
	return c.branchOnBitClear(0x80)
}

// -------- BBS0 - BBS7 --------

func (c *CPU6502) bbs0() (uint64, bool) {
	return c.branchOnBitSet(0x01)
}

func (c *CPU6502) bbs1() (uint64, bool) {
	return c.branchOnBitSet(0x02)
}

func (c *CPU6502) bbs2() (uint64, bool) {
	return c.branchOnBitSet(0x04)
}

func (c *CPU6502) bbs3() (uint64, bool) {
	return c.branchOnBitSet(0x08)
}

func (c *CPU6502) bbs4() (uint64, bool) {
	return c.branchOnBitSet(0x10)
}

func (c *CPU6502) bbs5() (uint64, bool) {
	return c.branchOnBitSet(0x20)
}

func (c *CPU6502) bbs6() (uint64, bool) {
	return c.branchOnBitSet(0x40)
}

func (c *CPU6502) bbs7() (uint64, bool) {
	return c.branchOnBitSet(0x80)
}

// -------- JSR--------

func (c *CPU6502) jsr() (uint64, bool) {
	addr := c.getAddrAbsolute()
	hiByte := uint8((c.PC & 0xFF00) >> 8)
	c.push(hiByte)
	loByte := uint8(c.PC & 0x00FF)
	c.push(loByte)
	c.PC = addr

	return 6, false
}

// -------- RTS--------

func (c *CPU6502) rts() (uint64, bool) {
	loByte := uint16(c.pop())
	hiByte := uint16(c.pop())
	addr := hiByte*256 + loByte + 1
	c.PC = addr

	return 6, false
}

// -------- JMP--------

func (c *CPU6502) jmp() (uint64, bool) {
	addr := c.getAddrAbsolute()
	c.PC = addr

	return 3, false
}

func (c *CPU6502) jmpIndirect6502() (uint64, bool) {
	addr := c.getAddrIndirectJmp6502()
	c.PC = addr

	return 5, false
}

func (c *CPU6502) jmpIndirect65C02() (uint64, bool) {
	addr := c.getAddrIndirect()
	c.PC = addr

	return 5, false
}

func (c *CPU6502) jmpIndexXIndirect() (uint64, bool) {
	addr := c.getAddrIdxIndirect65C02()
	c.PC = addr

	return 6, false
}
