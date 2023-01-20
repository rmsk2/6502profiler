package cpu

// -------- BPL --------

func (c *CPU6502) bpl() (uint64, bool) {
	if (c.Flags & Flag_N) != 0 {
		c.PC++
		return 2, false
	}

	branchAddress := c.getAddrRelative()
	// ToDo: Verify the assumption that the check for page crossing
	// is done relative to the first byte followng the branch
	// instruction
	additionalCycle := c.pageCrossCycles(c.PC+1, branchAddress)
	c.PC = branchAddress

	return 3 + additionalCycle, false
}
