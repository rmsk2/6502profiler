package cpu

// -------- STA --------

func (c *CPU6502) staAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.A)
	c.PC++

	return 4, false
}

func (c *CPU6502) staZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.A)
	c.PC++

	return 3, false
}

func (c *CPU6502) staAbsoluteY() (uint64, bool) {
	addr, _ := c.getAddrAbsoluteY()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 5, false
}

func (c *CPU6502) staAbsoluteX() (uint64, bool) {
	addr, _ := c.getAddrAbsoluteX()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 5, false
}

func (c *CPU6502) staZeroPageX() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageX(), c.A)
	c.PC++

	return 4, false
}

func (c *CPU6502) staIndirectY() (uint64, bool) {
	addr, _ := c.getAddrIndirectIdxY()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 6, false
}

func (c *CPU6502) staIndirect() (uint64, bool) {
	addr := c.getAddrZp65C02()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 5, false
}

func (c *CPU6502) staXIndirect() (uint64, bool) {
	addr := c.getAddrIdxIndirectX()
	c.Mem.Store(addr, c.A)
	c.PC++

	return 6, false
}

// -------- STZ --------

func (c *CPU6502) stzZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), 0)
	c.PC++

	return 4, false
}

func (c *CPU6502) stzZeroPageX() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageX(), 0)
	c.PC++

	return 5, false
}

func (c *CPU6502) stzAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), 0)
	c.PC++

	return 5, false
}

func (c *CPU6502) stzAbsoluteX() (uint64, bool) {
	addr, _ := c.getAddrAbsoluteX()
	c.Mem.Store(addr, 0)
	c.PC++

	return 6, false
}

// -------- STX --------

func (c *CPU6502) stxZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.X)
	c.PC++

	return 3, false
}

func (c *CPU6502) stxZeroPageY() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageY(), c.X)
	c.PC++

	return 4, false
}

func (c *CPU6502) stxAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.X)
	c.PC++

	return 4, false
}

// -------- STY --------

func (c *CPU6502) styZeroPage() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPage(), c.Y)
	c.PC++

	return 3, false
}

func (c *CPU6502) styZeroPageX() (uint64, bool) {
	c.Mem.Store(c.getAddrZeroPageX(), c.Y)
	c.PC++

	return 4, false
}

func (c *CPU6502) styAbsolute() (uint64, bool) {
	c.Mem.Store(c.getAddrAbsolute(), c.Y)
	c.PC++

	return 4, false
}

// -------- PHA --------

func (c *CPU6502) pha() (uint64, bool) {
	c.push(c.A)
	return 3, false
}

// -------- PLA --------

func (c *CPU6502) pla() (uint64, bool) {
	c.A = c.pop()
	c.nzFlags(c.A)
	return 4, false
}

// -------- PHX --------

func (c *CPU6502) phx() (uint64, bool) {
	c.push(c.X)
	return 3, false
}

// -------- PLX --------

func (c *CPU6502) plx() (uint64, bool) {
	c.X = c.pop()
	c.nzFlags(c.X)
	return 4, false
}

// -------- PHY --------

func (c *CPU6502) phy() (uint64, bool) {
	c.push(c.Y)
	return 3, false
}

// -------- PLY --------

func (c *CPU6502) ply() (uint64, bool) {
	c.Y = c.pop()
	c.nzFlags(c.Y)
	return 4, false
}

// -------- Flag stuff --------

func (c *CPU6502) clc() (uint64, bool) {
	c.Flags &= (^Flag_C)

	return 2, false
}

func (c *CPU6502) cli() (uint64, bool) {
	c.Flags &= (^Flag_I)

	return 2, false
}

func (c *CPU6502) clv() (uint64, bool) {
	c.Flags &= (^Flag_V)

	return 2, false
}

func (c *CPU6502) cld() (uint64, bool) {
	c.Flags &= (^Flag_D)

	return 2, false
}

func (c *CPU6502) sec() (uint64, bool) {
	c.Flags |= Flag_C

	return 2, false
}

func (c *CPU6502) sei() (uint64, bool) {
	c.Flags |= Flag_I

	return 2, false
}

func (c *CPU6502) sed() (uint64, bool) {
	c.Flags |= Flag_D

	return 2, false
}

// -------- PHP --------

func (c *CPU6502) php() (uint64, bool) {
	c.push(c.Flags)

	return 3, false
}

// -------- PLP --------

func (c *CPU6502) plp() (uint64, bool) {
	c.Flags = c.pop()

	return 4, false
}

// -------- TAX --------

func (c *CPU6502) tax() (uint64, bool) {
	c.nzFlags(c.A)
	c.X = c.A

	return 2, false
}

// -------- TXA --------

func (c *CPU6502) txa() (uint64, bool) {
	c.nzFlags(c.X)
	c.A = c.X

	return 2, false
}

// -------- TAY --------

func (c *CPU6502) tay() (uint64, bool) {
	c.nzFlags(c.A)
	c.Y = c.A

	return 2, false
}

// -------- TYA --------

func (c *CPU6502) tya() (uint64, bool) {
	c.nzFlags(c.Y)
	c.A = c.Y

	return 2, false
}

// -------- TXS --------

func (c *CPU6502) txs() (uint64, bool) {
	c.SP = c.X

	return 2, false
}

// -------- TSX --------

func (c *CPU6502) tsx() (uint64, bool) {
	c.nzFlags(c.SP)
	c.X = c.SP

	return 2, false
}
