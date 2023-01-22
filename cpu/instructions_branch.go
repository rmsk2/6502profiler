package cpu

// -------- BPL --------

func (c *CPU6502) bpl() (uint64, bool) {
	if (c.Flags & Flag_N) != 0 {
		c.PC++
		return 2, false
	}

	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}
