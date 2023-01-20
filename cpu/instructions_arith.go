package cpu

// -------- DEY --------

func (c *CPU6502) dey() (uint64, bool) {
	c.Y--
	c.nzFlags(c.Y)

	return 2, false
}
