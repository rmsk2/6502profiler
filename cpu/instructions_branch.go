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
