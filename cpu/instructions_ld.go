package cpu

// -------- LDX --------

func (c *CPU6502) ldxBase(value uint8) bool {
	c.X = value
	c.nzFlags(c.X)

	return false
}

func (c *CPU6502) ldxImmediate() (uint64, bool) {
	stop := c.ldxBase(c.Mem.Load(c.PC))
	c.PC++

	return 2, stop
}

func (c *CPU6502) ldxZeroPage() (uint64, bool) {
	stop := c.ldxBase(c.Mem.Load(c.getAddrZeroPage()))
	c.PC++

	return 3, stop
}

func (c *CPU6502) ldxZeroPageIdxY() (uint64, bool) {
	operandAddress := c.getAddrZeroPageY()
	stop := c.ldxBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldxAbsolute() (uint64, bool) {
	stop := c.ldxBase(c.Mem.Load(c.getAddrAbsolute()))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldxAbsoluteY() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteY()
	stop := c.ldxBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, stop
}

// -------- LDY --------

func (c *CPU6502) ldyBase(value uint8) bool {
	c.Y = value
	c.nzFlags(c.Y)

	return false
}

func (c *CPU6502) ldyImmediate() (uint64, bool) {
	stop := c.ldyBase(c.Mem.Load(c.PC))
	c.PC++

	return 2, stop
}

func (c *CPU6502) ldyZeroPage() (uint64, bool) {
	stop := c.ldyBase(c.Mem.Load(c.getAddrZeroPage()))
	c.PC++

	return 3, stop
}

func (c *CPU6502) ldyZeroPageIdxX() (uint64, bool) {
	operandAddress := c.getAddrZeroPageX()
	stop := c.ldyBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldyAbsolute() (uint64, bool) {
	stop := c.ldyBase(c.Mem.Load(c.getAddrAbsolute()))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldyAbsoluteX() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteX()
	stop := c.ldyBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, stop
}

// -------- LDA --------

func (c *CPU6502) ldaBase(value uint8) bool {
	c.A = value
	c.nzFlags(c.A)

	return false
}

func (c *CPU6502) ldaImmediate() (uint64, bool) {
	stop := c.ldaBase(c.Mem.Load(c.PC))
	c.PC++

	return 2, stop
}

func (c *CPU6502) ldaAbsolute() (uint64, bool) {
	stop := c.ldaBase(c.Mem.Load(c.getAddrAbsolute()))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldaZeroPage() (uint64, bool) {
	stop := c.ldaBase(c.Mem.Load(c.getAddrZeroPage()))
	c.PC++

	return 3, stop
}

func (c *CPU6502) ldaZeroPageIdxX() (uint64, bool) {
	operandAddress := c.getAddrZeroPageX()
	stop := c.ldaBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4, stop
}

func (c *CPU6502) ldaAbsoluteY() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteY()
	stop := c.ldaBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, stop
}

func (c *CPU6502) ldaAbsoluteX() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrAbsoluteX()
	stop := c.ldaBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, stop
}

func (c *CPU6502) ldaIndIdxY() (uint64, bool) {
	operandAddress, additionalCycle := c.getAddrIndirectIdxY()
	stop := c.ldaBase(c.Mem.Load(operandAddress))
	c.PC++

	return 4 + additionalCycle, stop
}

func (c *CPU6502) ldaIdxIndirectX() (uint64, bool) {
	operandAddress := c.getAddrIdxIndirectX()
	stop := c.ldaBase(c.Mem.Load(operandAddress))
	c.PC++

	return 6, stop
}
