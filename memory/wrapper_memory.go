package memory

type DataWriteFunc func(data uint8)

type MemWrapper interface {
	SetBaseMem(m Memory)
	Write(data uint8)
}

type WrappingMemory struct {
	mem                   Memory
	specialWriteAddresses map[uint16]DataWriteFunc
	ioMask                uint16
	wrappers              []MemWrapper
}

func NewMemWrapper(m Memory, mask uint16) *WrappingMemory {
	res := &WrappingMemory{
		mem:                   m,
		specialWriteAddresses: make(map[uint16]DataWriteFunc),
		ioMask:                mask & 0xFF00,
		wrappers:              []MemWrapper{},
	}

	return res
}

func (p *WrappingMemory) AddWrapper(addr uint16, wrapper MemWrapper) {
	p.specialWriteAddresses[addr] = wrapper.Write
	p.wrappers = append(p.wrappers, wrapper)
}

func (p *WrappingMemory) AddSpecialWriteAddress(addr uint16, f DataWriteFunc) {
	p.specialWriteAddresses[addr] = f
}

func (p *WrappingMemory) Load(address uint16) uint8 {
	return p.mem.Load(address)
}

func (p *WrappingMemory) Store(address uint16, b uint8) {
	if (address & 0xFF00) != p.ioMask {
		p.mem.Store(address, b)
		return
	}

	procFunc, ok := p.specialWriteAddresses[address]
	if !ok {
		p.mem.Store(address, b)
		return
	}

	procFunc(b)
}

func (p *WrappingMemory) GetStatistics(address uint16) uint64 {
	return p.mem.GetStatistics(address)
}

func (p *WrappingMemory) ClearStatistics() {
	p.mem.ClearStatistics()
}

func (p *WrappingMemory) TakeSnapshot() {
	p.mem.TakeSnapshot()
}

func (p *WrappingMemory) RestoreSnapshot() {
	p.mem.RestoreSnapshot()
}

// ------------------------------------------------------------------------------

type PlaceholderWrapper struct {
	Wrapper        *WrappingMemory
	f              DataWriteFunc
	baseMemory     Memory
	specialAddress uint16
}

func NewPlaceholderWrapper(m Memory, addr uint16) *PlaceholderWrapper {
	w := NewMemWrapper(m, addr&0xFF00)

	res := &PlaceholderWrapper{
		Wrapper:        w,
		f:              nil,
		baseMemory:     m,
		specialAddress: addr,
	}

	res.Wrapper.AddSpecialWriteAddress(res.specialAddress, res.Write)

	return res
}

func (p *PlaceholderWrapper) SetWriteFunc(fn DataWriteFunc) {
	p.f = fn
}

func (p *PlaceholderWrapper) Write(data uint8) {
	if p.f != nil {
		p.f(data)
	} else {
		p.baseMemory.Store(p.specialAddress, data)
	}
}
