package cpu

// -------- BPL --------

func (c *CPU6502) bpl() (uint64, bool) {
	if (c.Flags & Flag_N) != 0 {
		c.PC++
		return 2, false
	}

	// ToDo: Verify the assumption that the check for page crossing
	// is done relative to the first byte followng the branch
	// instruction
	branchAddress, additionalCycle := c.getAddrRelative()
	c.PC = branchAddress

	return 3 + additionalCycle, false
}
