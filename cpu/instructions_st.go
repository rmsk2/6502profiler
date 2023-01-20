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
	c.Mem.Store(c.getAddrAbsoluteY(), c.A)
	c.PC++

	return 5, false
}
