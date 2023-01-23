package cpu

// -------- BPL --------

func (c *CPU6502) bpl() (uint64, bool) {
	if (c.Flags & Flag_N) != 0 {
		c.PC++ // Skip when negative flag is set
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BMI --------

func (c *CPU6502) bmi() (uint64, bool) {
	if (c.Flags & Flag_N) == 0 {
		c.PC++ // Skip when negative flag is clear
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BEQ --------

func (c *CPU6502) beq() (uint64, bool) {
	if (c.Flags & Flag_Z) == 0 {
		c.PC++ // Skip when zero flag is clear
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BNE --------

func (c *CPU6502) bne() (uint64, bool) {
	if (c.Flags & Flag_Z) != 0 {
		c.PC++ // Skip when zero flag is set
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BCS--------

func (c *CPU6502) bcs() (uint64, bool) {
	if (c.Flags & Flag_C) == 0 {
		c.PC++ // Skip when carry flag is clear
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BCC --------

func (c *CPU6502) bcc() (uint64, bool) {
	if (c.Flags & Flag_C) != 0 {
		c.PC++ // Skip when carry flag is set
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BVS--------

func (c *CPU6502) bvs() (uint64, bool) {
	if (c.Flags & Flag_V) == 0 {
		c.PC++ // Skip when overflowflag is clear
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}

// -------- BVC --------

func (c *CPU6502) bvc() (uint64, bool) {
	if (c.Flags & Flag_V) != 0 {
		c.PC++ // Skip when overflow flag is set
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}
