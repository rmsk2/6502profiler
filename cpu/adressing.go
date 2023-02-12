package cpu

func (c *CPU6502) pageCrossCycles(addr1, addr2 uint16) uint64 {
	var additionalCycles uint64 = 0
	if (addr1 & 0xFF00) != (addr2 & 0xFF00) {
		additionalCycles = 1
	}

	return additionalCycles
}

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

func (c *CPU6502) getAddrAbsoluteY() (uint16, uint64) {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)
	var res = addr + uint16(c.Y)

	return res, c.pageCrossCycles(addr, res)
}

func (c *CPU6502) getAddrAbsoluteX() (uint16, uint64) {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)
	res := addr + uint16(c.X)

	return res, c.pageCrossCycles(res, addr)
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

// This was a bug of the original 6502 JMP(addr) implementation. When the address of an
// indirect JMP is the last byte on a page, i.e. 0xXYFF then the second byte is taken
// from 0xXY00.
func (c *CPU6502) getAddrIndirectJmp6502() uint16 {
	loByte := c.Mem.Load(c.PC)
	c.PC++
	var addr uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)
	loByte++
	var addr2 uint16 = uint16(c.Mem.Load(c.PC))*256 + uint16(loByte)

	return uint16(c.Mem.Load(addr2))*256 + uint16(c.Mem.Load(addr))
}

func (c *CPU6502) getAddrRelative() (uint16, uint64) {
	offset := int16(int8(c.Mem.Load(c.PC)))

	// Offsets are calculated relative to the byte following the instruction.
	// That's the reason for the plus one
	res := uint16(int16(c.PC+1) + offset)

	return res, c.pageCrossCycles(c.PC+1, res)
}

func (c *CPU6502) getAddrIndirectIdxY() (uint16, uint64) {
	zpAddrLo := c.Mem.Load(c.PC)
	zpAddrHi := zpAddrLo + 1 // Overflow is allowed

	var addr uint16 = uint16(c.Mem.Load(uint16(zpAddrHi)))*256 + uint16(c.Mem.Load(uint16(zpAddrLo)))
	var res = addr + uint16(c.Y)

	return res, c.pageCrossCycles(addr, res)
}

func (c *CPU6502) getAddrIdxIndirectX() uint16 {
	zpAddrLo := c.Mem.Load(c.PC) + c.X // Overflow is allowed
	zpAddrHi := zpAddrLo + 1           // Overflow is allowed

	return uint16(c.Mem.Load(uint16(zpAddrHi)))*256 + uint16(c.Mem.Load(uint16(zpAddrLo)))
}

func (c *CPU6502) getAddrZp65C02() uint16 {
	zpAddrLo := c.Mem.Load(c.PC)
	zpAddrHi := zpAddrLo + 1 // Overflow is allowed

	var addr uint16 = uint16(c.Mem.Load(uint16(zpAddrHi)))*256 + uint16(c.Mem.Load(uint16(zpAddrLo)))

	return addr
}

func (c *CPU6502) getAddrIdxIndirect65C02() uint16 {
	baseAddrLo := c.Mem.Load(c.PC)
	c.PC++

	baseAddr := uint16(c.Mem.Load(c.PC))*256 + uint16(baseAddrLo)
	baseAddr += uint16(c.X)

	return uint16(c.Mem.Load(baseAddr+1))*256 + uint16(c.Mem.Load(baseAddr))
}

func (c *CPU6502) getAddressesBitBranchRelative() (uint16, uint16, uint64) {
	zpAddr := c.getAddrZeroPage()
	c.PC++
	branchAddress, additionalCycle := c.getAddrRelative()

	return zpAddr, branchAddress, additionalCycle
}
