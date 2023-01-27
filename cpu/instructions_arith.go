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

// -------- add and subtract --------

func (c *CPU6502) addBaseBin(val1, val2 uint8) uint8 {
	v1 := uint16(val1)
	v2 := uint16(val2)

	var carry uint16 = 0

	if (c.Flags & Flag_C) != 0 {
		carry = 1
	}

	t := v1 + v2 + carry

	r := uint8(t & 0xFF)

	if r == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &= (^Flag_Z)
	}

	if (r & 0x80) != 0 {
		c.Flags |= Flag_N
	} else {
		c.Flags &= (^Flag_N)
	}

	if t >= 256 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	// Source http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	if (((val1 ^ r) & (val2 ^ r)) & 0x80) != 0 {
		c.Flags |= Flag_V
	} else {
		c.Flags &= (^Flag_V)
	}

	return r
}

// See http://www.6502.org/tutorials/decimal_mode.html Appendix A
func (cp *CPU6502) addBaseBcd6502(val1, val2 uint8) uint8 {
	a := uint16(val1)
	b := uint16(val2)
	var c uint16 = 0

	if (cp.Flags & Flag_C) != 0 {
		c = 1
	}

	// 1a. AL = (A & $0F) + (B & $0F) + C
	al := (a & 0x0F) + (b & 0x0F) + c
	// 1b. If AL >= $0A, then AL = ((AL + $06) & $0F) + $10
	if al >= 0x0A {
		al = ((al + 0x06) & 0x0F) + 0x10
	}
	// 1c. A = (A & $F0) + (B & $F0) + AL
	a = (a & 0xF0) + (b & 0xF0) + al
	// 1d. Note that A can be >= $100 at this point
	// 1e. If (A >= $A0), then A = A + $60
	if a >= 0xA0 {
		a = a + 0x60
	}
	// 1f. The accumulator result is the lower 8 bits of A
	r := uint8(a & 0xFF)

	if a >= 0x100 {
		cp.Flags |= Flag_C
	} else {
		cp.Flags &= (^Flag_C)
	}

	// Calculate setting of N, V and Z flags, which are invalid but if
	// we know how to do it in a precise manner why not implementing it
	// correctly
	cr := uint8(c)

	//2a. AL = (A & $0F) + (B & $0F) + C
	al2 := (val1 & 0x0F) + (val2 & 0x0F) + cr
	//2b. If AL >= $0A, then AL = ((AL + $06) & $0F) + $10
	if al2 >= 0x0A {
		al2 = ((al2 + 0x06) & 0x0F) + 0x10
	}
	//2c. A = (A & $F0) + (B & $F0) + AL, using signed (twos complement) arithmetic
	temp := (val1 & 0xF0) + (val2 & 0xF0) + al2

	//2e. The N flag result is 1 if bit 7 of A is 1, and is 0 if bit 7 if A is 0
	if (temp & 0x80) != 0 {
		cp.Flags |= Flag_N
	} else {
		cp.Flags &= (^Flag_N)
	}
	//2f. The V flag result is 1 if A < -128 or A > 127, and is 0 if -128 <= A <= 127
	v := int8(temp)
	if (v < -128) || (v > 127) {
		cp.Flags |= Flag_V
	} else {
		cp.Flags &= (^Flag_V)
	}

	// Zero flag is set as in binary mode
	a = uint16(val1)
	t2 := a + b + c
	if uint8(t2&0xFF) == 0 {
		cp.Flags |= Flag_Z
	} else {
		cp.Flags &= (^Flag_Z)
	}

	return r
}

func (c *CPU6502) subBaseBin(val1, val2 uint8) uint8 {
	// ... SBC simply takes the ones complement of the second value and then performs an ADC ...
	// Source http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	return c.addBaseBin(val1, val2^0xFF)
}

// See http://www.6502.org/tutorials/decimal_mode.html Appendix A
func (c *CPU6502) subBaseBcd6502(val1, val2 uint8) uint8 {
	return 0
}
