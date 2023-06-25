package memory

const X512K uint8 = 0
const X2048K uint8 = 1

type X16Memory struct {
	ramSelector *uint8
	romSelector *uint8

	baseMem      []byte
	bankedRAM8K  []byte
	bankedROM16K []byte

	baseMemSnapshot      []byte
	bankedRAM8KSnaphot   []byte
	bankedROM16KSnapshot []byte

	statBase      []uint64
	statBankedRam []uint64
	statBankedRom []uint64
}

func NewX16Memory(model uint8) *X16Memory {
	ramBlocks := 64

	if model == X2048K {
		ramBlocks = 256
	}

	lowMem := make([]byte, 40*1024)

	res := &X16Memory{
		ramSelector: &lowMem[0],
		romSelector: &lowMem[1],

		baseMem:      lowMem,
		bankedRAM8K:  make([]byte, ramBlocks*8192),
		bankedROM16K: make([]byte, 32*16384),

		baseMemSnapshot:      make([]byte, 40*1024),
		bankedRAM8KSnaphot:   make([]byte, ramBlocks*8192),
		bankedROM16KSnapshot: make([]byte, 32*16384),

		statBase:      make([]uint64, 40*1024),
		statBankedRam: make([]uint64, ramBlocks*8192),
		statBankedRom: make([]uint64, 32*16384),
	}

	*res.ramSelector = 1
	*res.romSelector = 0

	res.ClearStatistics()

	return res
}

func (x *X16Memory) ClearStatistics() {
	for i := 0; i < len(x.statBase); i++ {
		x.statBase[i] = 0
	}

	for i := 0; i < len(x.statBankedRam); i++ {
		x.statBankedRam[i] = 0
	}

	for i := 0; i < len(x.statBankedRom); i++ {
		x.statBankedRom[i] = 0
	}
}

func (x *X16Memory) calcIndex(address uint16) (*uint8, *uint64) {
	switch {
	case address < 0xA000:
		return &x.baseMem[address], &x.statBase[address]
	case address >= 0xC000:
		i := uint32(address-0xC000) + (uint32((*x.romSelector)&0x1f) * 16384)
		return &x.bankedROM16K[i], &x.statBankedRom[i]
	default:
		i := uint32(address-0xA000) + (uint32(*x.ramSelector) * 8192)
		return &x.bankedRAM8K[i], &x.statBankedRam[i]
	}
}

func (x *X16Memory) TakeSnapshot() {
	copy(x.baseMemSnapshot, x.baseMem)
	copy(x.bankedRAM8KSnaphot, x.bankedRAM8K)
	copy(x.bankedROM16KSnapshot, x.bankedROM16K)
}

func (x *X16Memory) RestoreSnapshot() {
	copy(x.baseMem, x.baseMemSnapshot)
	copy(x.bankedRAM8K, x.bankedRAM8KSnaphot)
	copy(x.bankedROM16K, x.bankedROM16KSnapshot)
}

func (x *X16Memory) Load(address uint16) uint8 {
	mem, stat := x.calcIndex(address)
	(*stat)++
	return *mem
}

func (x *X16Memory) Store(address uint16, b uint8) {
	mem, stat := x.calcIndex(address)
	(*stat)++
	*mem = b
}

func (x *X16Memory) GetStatistics(address uint16) uint64 {
	_, stat := x.calcIndex(address)
	return *stat

}

func (x *X16Memory) Close() {

}
