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
