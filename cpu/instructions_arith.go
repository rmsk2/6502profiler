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

func (c *CPU6502) addZeroPage() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrZeroPage())
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 3 + additionalCycles, false
}

func (c *CPU6502) subZeroPage() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrZeroPage())
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 3 + additionalCycles, false
}

func (c *CPU6502) addZeroPageX() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrZeroPageX())
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

func (c *CPU6502) subZeroPageX() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrZeroPageX())
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

func (c *CPU6502) addAbsolute() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrAbsolute())
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

func (c *CPU6502) subAbsolute() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrAbsolute())
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

func (c *CPU6502) addAbsoluteX() (uint64, bool) {
	addr, moreCycles := c.getAddrAbsoluteX()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) subAbsoluteX() (uint64, bool) {
	addr, moreCycles := c.getAddrAbsoluteX()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) addAbsoluteY() (uint64, bool) {
	addr, moreCycles := c.getAddrAbsoluteY()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) subAbsoluteY() (uint64, bool) {
	addr, moreCycles := c.getAddrAbsoluteY()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) addIndirect() (uint64, bool) {
	addr := c.getAddrZp65C02()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 5 + additionalCycles, false
}

func (c *CPU6502) subIndirect() (uint64, bool) {
	addr := c.getAddrZp65C02()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 5 + additionalCycles, false
}

func (c *CPU6502) addIndirectIdxY() (uint64, bool) {
	addr, moreCycles := c.getAddrIndirectIdxY()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) subIndirectIdxY() (uint64, bool) {
	addr, moreCycles := c.getAddrIndirectIdxY()
	operand := c.Mem.Load(addr)
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles + moreCycles, false
}

func (c *CPU6502) addIdxXIndirect() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrIdxIndirectX())
	res, additionalCycles := c.addBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

func (c *CPU6502) subIdxXIndirect() (uint64, bool) {
	operand := c.Mem.Load(c.getAddrIdxIndirectX())
	res, additionalCycles := c.subBase(c.A, operand)
	c.A = res
	c.PC++

	return 4 + additionalCycles, false
}

// -------- Logical operations --------

type LogicalOp func(a, b uint8) uint8

func Xor(a, b uint8) uint8 {
	return a ^ b
}

func And(a, b uint8) uint8 {
	return a & b
}

func Or(a, b uint8) uint8 {
	return a | b
}

func logicalImmediate(c *CPU6502, op LogicalOp) (uint64, bool) {
	operand := c.Mem.Load(c.PC)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 2, false
}

func logicalZeroPage(c *CPU6502, op LogicalOp) (uint64, bool) {
	address := c.getAddrZeroPage()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 3, false
}

func logicalZeroPageX(c *CPU6502, op LogicalOp) (uint64, bool) {
	address := c.getAddrZeroPageX()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 4, false
}

func logicalAbsolute(c *CPU6502, op LogicalOp) (uint64, bool) {
	address := c.getAddrAbsolute()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 4, false
}

func logicalAbsoluteX(c *CPU6502, op LogicalOp) (uint64, bool) {
	address, additionalCycles := c.getAddrAbsoluteX()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 4 + additionalCycles, false
}

func logicalAbsoluteY(c *CPU6502, op LogicalOp) (uint64, bool) {
	address, additionalCycles := c.getAddrAbsoluteY()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 4 + additionalCycles, false
}

func logicalIdxXIndirect(c *CPU6502, op LogicalOp) (uint64, bool) {
	address := c.getAddrIdxIndirectX()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 6, false
}

func logicalIndirectIdxY(c *CPU6502, op LogicalOp) (uint64, bool) {
	address, additionalCycles := c.getAddrIndirectIdxY()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 5 + additionalCycles, false
}

func logicalIndirect(c *CPU6502, op LogicalOp) (uint64, bool) {
	address := c.getAddrZp65C02()
	operand := c.Mem.Load(address)
	c.A = op(c.A, operand)
	c.nzFlags(c.A)
	c.PC++

	return 5, false
}

// -------- EOR --------

func (c *CPU6502) eorImmediate() (uint64, bool) {
	return logicalImmediate(c, Xor)
}

func (c *CPU6502) eorZeroPage() (uint64, bool) {
	return logicalZeroPage(c, Xor)
}

func (c *CPU6502) eorZeroPageX() (uint64, bool) {
	return logicalZeroPageX(c, Xor)
}

func (c *CPU6502) eorAbsolute() (uint64, bool) {
	return logicalAbsolute(c, Xor)
}

func (c *CPU6502) eorAbsoluteX() (uint64, bool) {
	return logicalAbsoluteX(c, Xor)
}

func (c *CPU6502) eorAbsoluteY() (uint64, bool) {
	return logicalAbsoluteY(c, Xor)
}

func (c *CPU6502) eorIdxIndirect() (uint64, bool) {
	return logicalIdxXIndirect(c, Xor)
}

func (c *CPU6502) eorIndirectIdxY() (uint64, bool) {
	return logicalIndirectIdxY(c, Xor)
}

func (c *CPU6502) eorIndirect() (uint64, bool) {
	return logicalIndirect(c, Xor)
}

// -------- ORA --------

func (c *CPU6502) oraImmediate() (uint64, bool) {
	return logicalImmediate(c, Or)
}

func (c *CPU6502) oraZeroPage() (uint64, bool) {
	return logicalZeroPage(c, Or)
}

func (c *CPU6502) oraZeroPageX() (uint64, bool) {
	return logicalZeroPageX(c, Or)
}

func (c *CPU6502) oraAbsolute() (uint64, bool) {
	return logicalAbsolute(c, Or)
}

func (c *CPU6502) oraAbsoluteX() (uint64, bool) {
	return logicalAbsoluteX(c, Or)
}

func (c *CPU6502) oraAbsoluteY() (uint64, bool) {
	return logicalAbsoluteY(c, Or)
}

func (c *CPU6502) oraIdxIndirect() (uint64, bool) {
	return logicalIdxXIndirect(c, Or)
}

func (c *CPU6502) oraIndirectIdxY() (uint64, bool) {
	return logicalIndirectIdxY(c, Or)
}

func (c *CPU6502) oraIndirect() (uint64, bool) {
	return logicalIndirect(c, Or)
}

// -------- AND --------

func (c *CPU6502) andImmediate() (uint64, bool) {
	return logicalImmediate(c, And)
}

func (c *CPU6502) andZeroPage() (uint64, bool) {
	return logicalZeroPage(c, And)
}

func (c *CPU6502) andZeroPageX() (uint64, bool) {
	return logicalZeroPageX(c, And)
}

func (c *CPU6502) andAbsolute() (uint64, bool) {
	return logicalAbsolute(c, And)
}

func (c *CPU6502) andAbsoluteX() (uint64, bool) {
	return logicalAbsoluteX(c, And)
}

func (c *CPU6502) andAbsoluteY() (uint64, bool) {
	return logicalAbsoluteY(c, And)
}

func (c *CPU6502) andIdxIndirect() (uint64, bool) {
	return logicalIdxXIndirect(c, And)
}

func (c *CPU6502) andIndirectIdxY() (uint64, bool) {
	return logicalIndirectIdxY(c, And)
}

func (c *CPU6502) andIndirect() (uint64, bool) {
	return logicalIndirect(c, And)
}

// -------- Modifier operations --------

type ModifierOp func(c *CPU6502, a uint8) uint8

func Rol(c *CPU6502, a uint8) uint8 {
	var val uint8 = 0

	if (c.Flags & Flag_C) != 0 {
		val = 0x01
	}

	if (a & 0x80) != 0 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return (a << 1) | val
}

func Ror(c *CPU6502, a uint8) uint8 {
	var val uint8 = 0

	if (c.Flags & Flag_C) != 0 {
		val = 0x80
	}

	if (a & 1) != 0 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return (a >> 1) | val
}

func Lsr(c *CPU6502, a uint8) uint8 {
	if (a & 1) != 0 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return a >> 1
}

func Asl(c *CPU6502, a uint8) uint8 {
	if (a & 0x80) != 0 {
		c.Flags |= Flag_C
	} else {
		c.Flags &= (^Flag_C)
	}

	return a << 1
}

func Inc(c *CPU6502, a uint8) uint8 {
	return a + 1
}

func Dec(c *CPU6502, a uint8) uint8 {
	return a + 0xFF // 0xFF == -1
}

func (c *CPU6502) modImplied(modifier ModifierOp) (uint64, bool) {
	c.A = modifier(c, c.A)
	c.nzFlags(c.A)

	return 2, false
}

func (c *CPU6502) modZeroPage(modifier ModifierOp) (uint64, bool) {
	operAddr := c.getAddrZeroPage()
	oper := c.Mem.Load(operAddr)
	res := modifier(c, oper)
	c.Mem.Store(operAddr, res)
	c.nzFlags(res)
	c.PC++

	return 5, false
}

func (c *CPU6502) modZeroPageX(modifier ModifierOp) (uint64, bool) {
	operAddr := c.getAddrZeroPageX()
	oper := c.Mem.Load(operAddr)
	res := modifier(c, oper)
	c.Mem.Store(operAddr, res)
	c.nzFlags(res)
	c.PC++

	return 6, false
}

func (c *CPU6502) modAbsolute(modifier ModifierOp) (uint64, bool) {
	operAddr := c.getAddrAbsolute()
	oper := c.Mem.Load(operAddr)
	res := modifier(c, oper)
	c.Mem.Store(operAddr, res)
	c.nzFlags(res)
	c.PC++

	return 6, false
}

func (c *CPU6502) modAbsoluteX(modifier ModifierOp) (uint64, bool) {
	operAddr, _ := c.getAddrAbsoluteX()
	oper := c.Mem.Load(operAddr)
	res := modifier(c, oper)
	c.Mem.Store(operAddr, res)
	c.nzFlags(res)
	c.PC++

	return 7, false
}

// -------- INC --------

func (c *CPU6502) inc65C02() (uint64, bool) {
	return c.modImplied(Inc)
}

func (c *CPU6502) incZeroPage() (uint64, bool) {
	return c.modZeroPage(Inc)
}

func (c *CPU6502) incZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Inc)
}

func (c *CPU6502) incAbsolute() (uint64, bool) {
	return c.modAbsolute(Inc)
}

func (c *CPU6502) incAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Inc)
}

// -------- DEC --------

func (c *CPU6502) dec65C02() (uint64, bool) {
	return c.modImplied(Dec)
}

func (c *CPU6502) decZeroPage() (uint64, bool) {
	return c.modZeroPage(Dec)
}

func (c *CPU6502) decZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Dec)
}

func (c *CPU6502) decAbsolute() (uint64, bool) {
	return c.modAbsolute(Dec)
}

func (c *CPU6502) decAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Dec)
}

// -------- ASL --------

func (c *CPU6502) asl() (uint64, bool) {
	return c.modImplied(Asl)
}

func (c *CPU6502) aslZeroPage() (uint64, bool) {
	return c.modZeroPage(Asl)
}

func (c *CPU6502) aslZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Asl)
}

func (c *CPU6502) aslAbsolute() (uint64, bool) {
	return c.modAbsolute(Asl)
}

func (c *CPU6502) aslAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Asl)
}

// -------- LSR --------

func (c *CPU6502) lsr() (uint64, bool) {
	return c.modImplied(Lsr)
}

func (c *CPU6502) lsrZeroPage() (uint64, bool) {
	return c.modZeroPage(Lsr)
}

func (c *CPU6502) lsrZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Lsr)
}

func (c *CPU6502) lsrAbsolute() (uint64, bool) {
	return c.modAbsolute(Lsr)
}

func (c *CPU6502) lsrAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Lsr)
}

// -------- ROL --------

func (c *CPU6502) rol() (uint64, bool) {
	return c.modImplied(Rol)
}

func (c *CPU6502) rolZeroPage() (uint64, bool) {
	return c.modZeroPage(Rol)
}

func (c *CPU6502) rolZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Rol)
}

func (c *CPU6502) rolAbsolute() (uint64, bool) {
	return c.modAbsolute(Rol)
}

func (c *CPU6502) rolAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Rol)
}

// -------- ROR --------

func (c *CPU6502) ror() (uint64, bool) {
	return c.modImplied(Ror)
}

func (c *CPU6502) rorZeroPage() (uint64, bool) {
	return c.modZeroPage(Ror)
}

func (c *CPU6502) rorZeroPageX() (uint64, bool) {
	return c.modZeroPageX(Ror)
}

func (c *CPU6502) rorAbsolute() (uint64, bool) {
	return c.modAbsolute(Ror)
}

func (c *CPU6502) rorAbsoluteX() (uint64, bool) {
	return c.modAbsoluteX(Ror)
}

// -------- BIT --------

func (c *CPU6502) bitBase(val uint8) {
	if (c.A & val) == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &= (^Flag_Z)
	}

	if (val & 0x80) != 0 {
		c.Flags |= Flag_N
	} else {
		c.Flags &= (^Flag_N)
	}

	if (val & 0x40) != 0 {
		c.Flags |= Flag_V
	} else {
		c.Flags &= (^Flag_V)
	}
}

func (c *CPU6502) bitZeroPage() (uint64, bool) {
	oper := c.Mem.Load(c.getAddrZeroPage())
	c.bitBase(oper)
	c.PC++

	return 3, false
}

func (c *CPU6502) bitZeroPageX() (uint64, bool) {
	oper := c.Mem.Load(c.getAddrZeroPageX())
	c.bitBase(oper)
	c.PC++

	return 4, false
}

func (c *CPU6502) bitImmediate() (uint64, bool) {
	oper := c.Mem.Load(c.PC)
	c.bitBase(oper)
	c.PC++

	return 2, false
}

func (c *CPU6502) bitAbsolute() (uint64, bool) {
	oper := c.Mem.Load(c.getAddrAbsolute())
	c.bitBase(oper)
	c.PC++

	return 4, false
}

func (c *CPU6502) bitAbsoluteX() (uint64, bool) {
	addr, additionalCycles := c.getAddrAbsoluteX()
	oper := c.Mem.Load(addr)
	c.bitBase(oper)
	c.PC++

	return 4 + additionalCycles, false
}

// -------- TRB --------

func (c *CPU6502) trbBase(val uint8) uint8 {
	if (c.A & val) == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &= (^Flag_Z)
	}

	return (^c.A) & val
}

func (c *CPU6502) trbZeroPage() (uint64, bool) {
	addr := c.getAddrZeroPage()

	oper := c.Mem.Load(addr)
	res := c.trbBase(oper)

	c.Mem.Store(addr, res)
	c.PC++

	return 5, false
}

func (c *CPU6502) trbAbsolute() (uint64, bool) {
	addr := c.getAddrAbsolute()

	oper := c.Mem.Load(addr)
	res := c.trbBase(oper)

	c.Mem.Store(addr, res)
	c.PC++

	return 6, false
}

// -------- TSB --------

func (c *CPU6502) tsbBase(val uint8) uint8 {
	if (c.A & val) == 0 {
		c.Flags |= Flag_Z
	} else {
		c.Flags &= (^Flag_Z)
	}

	return c.A | val
}

func (c *CPU6502) tsbZeroPage() (uint64, bool) {
	addr := c.getAddrZeroPage()

	oper := c.Mem.Load(addr)
	res := c.tsbBase(oper)

	c.Mem.Store(addr, res)
	c.PC++

	return 5, false
}

func (c *CPU6502) tsbAbsolute() (uint64, bool) {
	addr := c.getAddrAbsolute()

	oper := c.Mem.Load(addr)
	res := c.tsbBase(oper)

	c.Mem.Store(addr, res)
	c.PC++

	return 6, false
}

// -------- RMB0-7 --------

func (c *CPU6502) rmbBase(bit uint8) (uint64, bool) {
	addr := c.getAddrZeroPage()

	oper := c.Mem.Load(addr)
	res := oper & (bit ^ 0xFF)

	c.Mem.Store(addr, res)
	c.PC++

	return 5, false
}

func (c *CPU6502) rmb0() (uint64, bool) {
	return c.rmbBase(0x01)
}

func (c *CPU6502) rmb1() (uint64, bool) {
	return c.rmbBase(0x02)
}

func (c *CPU6502) rmb2() (uint64, bool) {
	return c.rmbBase(0x04)
}

func (c *CPU6502) rmb3() (uint64, bool) {
	return c.rmbBase(0x08)
}

func (c *CPU6502) rmb4() (uint64, bool) {
	return c.rmbBase(0x10)
}

func (c *CPU6502) rmb5() (uint64, bool) {
	return c.rmbBase(0x20)
}

func (c *CPU6502) rmb6() (uint64, bool) {
	return c.rmbBase(0x40)
}

func (c *CPU6502) rmb7() (uint64, bool) {
	return c.rmbBase(0x80)
}

// -------- SMB0-7 --------

func (c *CPU6502) smbBase(bit uint8) (uint64, bool) {
	addr := c.getAddrZeroPage()

	oper := c.Mem.Load(addr)
	res := oper | bit

	c.Mem.Store(addr, res)
	c.PC++

	return 5, false
}

func (c *CPU6502) smb0() (uint64, bool) {
	return c.smbBase(0x01)
}

func (c *CPU6502) smb1() (uint64, bool) {
	return c.smbBase(0x02)
}

func (c *CPU6502) smb2() (uint64, bool) {
	return c.smbBase(0x04)
}

func (c *CPU6502) smb3() (uint64, bool) {
	return c.smbBase(0x08)
}

func (c *CPU6502) smb4() (uint64, bool) {
	return c.smbBase(0x10)
}

func (c *CPU6502) smb5() (uint64, bool) {
	return c.smbBase(0x20)
}

func (c *CPU6502) smb6() (uint64, bool) {
	return c.smbBase(0x40)
}

func (c *CPU6502) smb7() (uint64, bool) {
	return c.smbBase(0x80)
}
