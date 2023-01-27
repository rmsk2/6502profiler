package cpu

// -------- STA --------

func (c *CPU6502) staAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.A)
	c.PC++

	return 4, false
}

func (c *CPU6502) staZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.A)
	c.PC++

	return 3, false
}

func (c *CPU6502) staAbsoluteY() (uint64, bool) {
	addr, _ := c.getAddrAbsoluteY()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 5, false
}

func (c *CPU6502) staAbsoluteX() (uint64, bool) {
	addr, _ := c.getAddrAbsoluteX()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 5, false
}

func (c *CPU6502) staZeroPageX() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageX(), c.A)
	c.PC++

	return 4, false
}

func (c *CPU6502) staIndirectY() (uint64, bool) {
	addr, _ := c.getAddrIndirectIdxY()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 6, false
}

func (c *CPU6502) staXIndirect() (uint64, bool) {
	addr := c.getAddrIdxIndirectX()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 6, false
}

// -------- STX --------

func (c *CPU6502) stxZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.X)
	c.PC++

	return 3, false
}

func (c *CPU6502) stxZeroPageY() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageY(), c.X)
	c.PC++

	return 4, false
}

func (c *CPU6502) stxAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.X)
	c.PC++

	return 4, false
}

// -------- STY --------

func (c *CPU6502) styZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.Y)
	c.PC++

	return 3, false
}

func (c *CPU6502) styZeroPageX() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageX(), c.Y)
	c.PC++

	return 4, false
}

func (c *CPU6502) styAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.Y)
	c.PC++

	return 4, false
}

// -------- PHA --------

func (c *CPU6502) pha() (uint64, bool) {
	c.push(c.A)
	return 3, false
}

// -------- PLA --------

func (c *CPU6502) pla() (uint64, bool) {
	c.A = c.pop()
	c.nzFlags(c.A)
	return 4, false
}
