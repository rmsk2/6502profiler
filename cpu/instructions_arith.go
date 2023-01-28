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

	c.nzFlags(r)

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

func (c *CPU6502) prepBCD(val1, val2 uint8) (uint8, uint8, uint8) {
	b1 := c.fromBCD(val1)
	b2 := c.fromBCD(val2)

	var carry uint8 = 0
	if (c.Flags & Flag_C) != 0 {
		carry = 1
	}

	return b1, b2, carry
}

func (c *CPU6502) fromBCD(in uint8) uint8 {
	loNibble := in & 0x0F
	hiNibble := (in & 0xF0) >> 4

	if (loNibble > 9) || (hiNibble > 9) {
		panic("Invalid BCD data")
	}

	return (hiNibble * 10) + loNibble
}

func (c *CPU6502) toBCD(res uint8) uint8 {
	res = res % 100
	loDigit := res % 10
	hiDigit := res / 10

	return (hiDigit << 4) | loDigit
}

// See http://www.6502.org/tutorials/decimal_mode.html
// We implement a logically correct version of BCD. This routine does not simulate undocumented
// behaviour. We therefore panic if the number is not valid BCD.
// On a 6502 the only documented behaviour is for
// setting the carry flag. On the 65C02 the N and Z flags are also valid. This routine implements
// the 65C02 behaviour as it is a superset of the 6502 behaviour.
func (c *CPU6502) addBaseBcd6502(val1, val2 uint8) uint8 {
	b1, b2, carry := c.prepBCD(val1, val2)

	cr := b1 + b2 + carry
	res := c.toBCD(cr)

	c.nzFlags(res)

	if cr >= 100 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return res
}

// See http://www.6502.org/tutorials/decimal_mode.html
// We implement a logically correct version of BCD. This routine does not simulate undocumented
// behaviour. We therefore panic if the number is not valid BCD.
// On a 6502 the only documented behaviour is for
// setting the carry flag. On the 65C02 the N and Z flags are also valid. This routine implements
// the 65C02 behaviour as it is a superset of the 6502 behaviour.
func (c *CPU6502) subBaseBcd(val1, val2 uint8) uint8 {
	var res uint8
	t1, t2, carry := c.prepBCD(val1, val2)

	b1 := int16(t1)
	b2 := int16(t2)
	cr := int16(1 - carry)

	temp := b1 - b2 - cr

	if temp < 0 {
		c.Flags &= (^Flag_C)
		res = c.toBCD(uint8(temp + 100))
		c.nzFlags(res)
	} else {
		c.Flags |= Flag_C
		res = c.toBCD(uint8(temp))
		c.nzFlags(res)
	}

	return res
}

func (c *CPU6502) addBase(val1, val2 uint8) (uint8, uint64) {
	var res uint8
	var additionalCycles uint64 = 0

	if (c.model == Model65C02) && ((c.Flags & Flag_D) != 0) {
		additionalCycles++
	}

	if (c.Flags & Flag_D) == 0 {
		res = c.addBaseBin(val1, val2)
	} else {
		res = c.addBaseBcd6502(val1, val2)
	}

	return res, additionalCycles
}

func (c *CPU6502) subBase(val1, val2 uint8) (uint8, uint64) {
	var res uint8
	var additionalCycles uint64 = 0

	if (c.model == Model65C02) && ((c.Flags & Flag_D) != 0) {
		additionalCycles++
	}

	if (c.Flags & Flag_D) == 0 {
		res = c.subBaseBin(val1, val2)
	} else {
		res = c.subBaseBcd(val1, val2)
	}

	return res, additionalCycles
}

func (c *CPU6502) addImmediate() (uint64, bool) {
	operand := c.Mem.Load(c.PC)
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 2 + additionalCycles, false
}

func (c *CPU6502) subImmediate() (uint64, bool) {
	operand := c.Mem.Load(c.PC)
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 2 + additionalCycles, false
}
