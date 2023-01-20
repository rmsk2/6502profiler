package cpu

// -------- STA --------

func (c *CPU6502) staAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.A)
	c.PC++

	return 4, false
}
