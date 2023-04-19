package memory

type F256UnsignedCoproc struct {
	baseAdress uint16
	mem        Memory
}

func NewUnsignedCoproc(baseMem Memory, base uint16) *F256UnsignedCoproc {
	return &F256UnsignedCoproc{
		baseAdress: base,
		mem:        baseMem,
	}
}

func (f *F256UnsignedCoproc) RegisterUmul(wrapper *WrappingMemory) {
	var i uint16

	for i = 0; i < 4; i++ {
		var targetAddress uint16 = f.baseAdress + i
		wrapper.AddSpecialWriteAddress(targetAddress, func(data uint8) { f.WriteUMul(targetAddress, data) })
	}
}

func (f *F256UnsignedCoproc) RegisterUdiv(wrapper *WrappingMemory) {
	var i uint16

	for i = 0; i < 4; i++ {
		var targetAddress uint16 = f.baseAdress + 8 + i
		wrapper.AddSpecialWriteAddress(targetAddress, func(data uint8) { f.WriteUDiv(targetAddress, data) })
	}
}

func (f *F256UnsignedCoproc) WriteUMul(address uint16, data uint8) {
	f.mem.Store(address, data)
	oper1 := uint(f.mem.Load(f.baseAdress+1))*256 + uint(f.mem.Load(f.baseAdress))
	oper2 := uint(f.mem.Load(f.baseAdress+3))*256 + uint(f.mem.Load(f.baseAdress+2))

	res := oper1 * oper2
	var i uint16

	for i = 0; i < 4; i++ {
		f.mem.Store(f.baseAdress+4+i, uint8(res&0xFF))
		res >>= 8
	}
}

func (f *F256UnsignedCoproc) WriteUDiv(address uint16, data uint8) {
	f.mem.Store(address, data)
	oper1 := uint(f.mem.Load(f.baseAdress+9))*256 + uint(f.mem.Load(f.baseAdress+8))
	oper2 := uint(f.mem.Load(f.baseAdress+11))*256 + uint(f.mem.Load(f.baseAdress+10))

	resDiv := oper2 / oper1
	resMod := oper2 % oper1

	f.mem.Store(f.baseAdress+0x14, uint8(resDiv&0xFF))
	f.mem.Store(f.baseAdress+0x15, uint8((resDiv&0xFF00)>>8))
	f.mem.Store(f.baseAdress+0x16, uint8(resMod&0xFF))
	f.mem.Store(f.baseAdress+0x17, uint8((resMod&0xFF00)>>8))
}
