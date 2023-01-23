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
