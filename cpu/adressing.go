package cpu

// -------- Addressing modes --------

func (c *CPU6502) getAddrAbsolute() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return addr
}

func (c *CPU6502) getAddrZeroPage() uint16 {
	loByte := c.Mem.Load(c.PC)

	return uint16(loByte)
}

func (c *CPU6502) getAddrAbsoluteY() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return addr + uint16(c.Y)
}

func (c *CPU6502) getAddrAbsoluteX() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return addr + uint16(c.X)
}

func (c *CPU6502) getAddrZeroPageY() uint16 {
	loByte := c.Mem.Load(c.PC)
	var zpAddr uint8 = loByte + c.Y // Allow possible overflow

	return uint16(zpAddr)
}

func (c *CPU6502) getAddrZeroPageX() uint16 {
	loByte := c.Mem.Load(c.PC)
	var zpAddr uint8 = loByte + c.X // Allow possible overflow

	return uint16(zpAddr)
}

func (c *CPU6502) getAddrIndirect() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return uint16(c.Mem.Load(addr+1))*256 + uint16(c.Mem.Load(addr))
}

func (c *CPU6502) getAddrRelative() uint16 {
	offset := int16(int8(c.Mem.Load(c.PC)))

	// Offsets are calculated relative to the byte following the instruction.
	// That's the reason for the plus one
	return uint16(int16(c.PC+1) + offset)
}

func (c *CPU6502) getAddrIndirectIdxY() uint16 {
	zpAddrLo := c.Mem.Load(c.PC)
	zpAddrHi := zpAddrLo + 1 // Overflow is allowed

	var addr uint16 = uint16(c.Mem.Load(uint16(zpAddrHi)))*256 + uint16(c.Mem.Load(uint16(zpAddrLo)))

	return addr + uint16(c.Y)
}

func (c *CPU6502) getAddrIdxIndirectX() uint16 {
	zpAddrLo := c.Mem.Load(c.PC) + c.X // Overflow is allowed
	zpAddrHi := zpAddrLo + 1           // Overflow is allowed

	return uint16(c.Mem.Load(uint16(zpAddrHi)))*256 + uint16(c.Mem.Load(uint16(zpAddrLo)))
}
