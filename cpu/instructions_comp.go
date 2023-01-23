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
