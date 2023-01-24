package cpu

// -------- DEY --------

func (c *CPU6502) dey() (uint64, bool) {
	c.Y--
	c.nzFlags(c.Y)

	return 2, false
}

// -------- INY --------

func (c *CPU6502) iny() (uint64, bool) {
	c.Y++
	c.nzFlags(c.Y)

	return 2, false
}

// -------- DEX --------

func (c *CPU6502) dex() (uint64, bool) {
	c.X--
	c.nzFlags(c.X)

	return 2, false
}

// -------- INX --------

func (c *CPU6502) inx() (uint64, bool) {
	c.X++
	c.nzFlags(c.X)

	return 2, false
}
