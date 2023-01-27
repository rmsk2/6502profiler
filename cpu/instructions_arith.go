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

func (c *CPU6502) subBaseBin(val1, val2 uint8) uint8 {
	// ... SBC simply takes the ones complement of the second value and then performs an ADC ...
	// Source http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	return c.addBaseBin(val1, val2^0xFF)
}

func (c *CPU6502) fromBCD(val1, val2 uint8) (uint8, uint8, uint8, uint8, uint8) {
	loNibble1 := uint8(val1 & 0x0F)
	hiNibble1 := uint8((val1 & 0xF0) >> 4)

	loNibble2 := uint8(val2 & 0x0F)
	hiNibble2 := uint8((val2 & 0xF0) >> 4)

	var carry uint8 = 0
	if (c.Flags & Flag_C) != 0 {
		carry = 1
	}

	if (loNibble1 > 9) || (hiNibble1 > 9) || (loNibble2 > 9) || (hiNibble2 > 9) {
		panic("Invalid BCD data")
	}

	return loNibble1, hiNibble1, loNibble2, hiNibble2, carry
}

func (c *CPU6502) toBCD(res uint8) uint8 {
	res = res % 100
	loDigit := res % 10
	hiDigit := res / 10

	return (hiDigit >> 4) | loDigit
}

// See http://www.6502.org/tutorials/decimal_mode.html
// We implement a logically correct version of BCD. This routine does not simulate undocumented
// behaviour. We therefore panic if the number is not valid BCD and we do not change any
// flags other than the carry flag. On a 6502 the only documented behaviour is for
// setting the carry flag.
func (c *CPU6502) addBaseBcd6502(val1, val2 uint8) uint8 {
	loNibble1, hiNibble1, loNibble2, hiNibble2, carry := c.fromBCD(val1, val2)

	cr := hiNibble1*10 + loNibble1 + hiNibble2*10 + loNibble2 + carry

	if (cr / 100) > 0 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return c.toBCD(cr % 100)
}

// See http://www.6502.org/tutorials/decimal_mode.html
// We implement a logically correct version of BCD. This routine does not simulate undocumented
// behaviour. We therefore panic if the number is not valid BCD and we do not change any
// flags other than the carry flag. On a 6502 the only documented behaviour is for
// setting the carry flag.
func (c *CPU6502) subBaseBcd6502(val1, val2 uint8) uint8 {
	loNibble1, hiNibble1, loNibble2, hiNibble2, carry := c.fromBCD(val1, val2)

	b1 := int16(hiNibble1*10 + loNibble1)
	b2 := int16(hiNibble2*10 + loNibble2)
	cr := int16(1 - carry)

	temp := b1 - b2 - cr

	if temp < 0 {
		c.Flags &= (^Flag_C)
		return c.toBCD(uint8(temp + 100))
	} else {
		c.Flags |= Flag_C
		return c.toBCD(uint8(temp))
	}
}
