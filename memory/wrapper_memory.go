package memory

type DataWriteFunc func(data uint8)

type WrappingMemory struct {
	mem              Memory
	specialAddresses map[uint16]DataWriteFunc
}

func NewMemWrapper(m Memory) *WrappingMemory {
	res := &WrappingMemory{
		mem:              m,
		specialAddresses: make(map[uint16]DataWriteFunc),
	}

	return res
}

func (p *WrappingMemory) AddSpecialWriteAddress(addr uint16, f DataWriteFunc) {
	p.specialAddresses[addr] = f
}

func (p *WrappingMemory) Load(address uint16) uint8 {
	return p.mem.Load(address)
}

func (p *WrappingMemory) Store(address uint16, b uint8) {
	procFunc, ok := p.specialAddresses[address]
	if !ok {
		p.mem.Store(address, b)
		return
	}

	procFunc(b)
}

func (p *WrappingMemory) GetStatistics() []BankAccessStatistics {
	return p.mem.GetStatistics()
}
