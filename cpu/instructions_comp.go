package cpu

func (c *CPU6502) cmpBase(val1, val2 uint8) {
	if val1 == val2 {
		c.Flags |= Flag_Z    // Zero is set
		c.Flags |= Flag_C    // Carry is set
		c.Flags &= (^Flag_N) // Negative is clear
		return
	}

	t := val1 - val2
	if (t & 0x80) != 0 {
		c.Flags |= Flag_N // Negative is set
	} else {
		c.Flags &= (^Flag_N) // Negative is clear
	}

	if val1 > val2 {
		c.Flags &= (^Flag_Z) // Zero is clear
		c.Flags |= Flag_C    // Carry is set
		return
	}

	// val1 < val2
	c.Flags &= (^Flag_Z) // Zero is clear
	c.Flags &= (^Flag_C) // Carry is clear
}

// -------- CPY --------

func (c *CPU6502) cpyImmediate() (uint64, bool) {
	c.cmpBase(c.Y, c.Mem.Load(c.PC))
	c.PC++

	return 2, false
}

func (c *CPU6502) cpyZeroPage() (uint64, bool) {
	addr := c.getAddrZeroPage()
	c.cmpBase(c.Y, c.Mem.Load(addr))
	c.PC++

	return 3, false
}

func (c *CPU6502) cpyAbsolute() (uint64, bool) {
	addr := c.getAddrAbsolute()
	c.cmpBase(c.Y, c.Mem.Load(addr))
	c.PC++

	return 4, false
}

// -------- CPX --------

func (c *CPU6502) cpxImmediate() (uint64, bool) {
	c.cmpBase(c.X, c.Mem.Load(c.PC))
	c.PC++

	return 2, false
}

func (c *CPU6502) cpxZeroPage() (uint64, bool) {
	addr := c.getAddrZeroPage()
	c.cmpBase(c.X, c.Mem.Load(addr))
	c.PC++

	return 3, false
}

func (c *CPU6502) cpxAbsolute() (uint64, bool) {
	addr := c.getAddrAbsolute()
	c.cmpBase(c.X, c.Mem.Load(addr))
	c.PC++

	return 4, false
}

// -------- CMP --------

func (c *CPU6502) cmpImmediate() (uint64, bool) {
	c.cmpBase(c.A, c.Mem.Load(c.PC))
	c.PC++

	return 2, false
}

func (c *CPU6502) cmpZeroPage() (uint64, bool) {
	operandAddress := c.getAddrZeroPage()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 3, false
}

func (c *CPU6502) cmpZeroPageX() (uint64, bool) {
	operandAddress := c.getAddrZeroPageX()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 4, false
}

func (c *CPU6502) cmpAbsolute() (uint64, bool) {
	operandAddress := c.getAddrAbsolute()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 4, false
}

func (c *CPU6502) cmpAbsoluteX() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteX()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, false
}

func (c *CPU6502) cmpAbsoluteY() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteY()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, false
}

func (c *CPU6502) cmpIdxXIndirect() (uint64, bool) {
	operandAddress := c.getAddrIdxIndirectX()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 6, false
}

func (c *CPU6502) cmpIndIdxY() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrIndirectIdxY()
	c.cmpBase(c.A, c.Mem.Load(operandAddress))
	c.PC++

	return 5 + additionalCycle, false
}
