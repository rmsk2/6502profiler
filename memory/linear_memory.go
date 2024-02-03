package memory

type LinearMemory struct {
	memory         []byte
	accessCount    []uint64
	memorySnapshot []byte
}

func NewLinearMemory(size uint32) *LinearMemory {
	if size > 65536 {
		size = 65536
	}

	res := &LinearMemory{
		memory:         make([]byte, size),
		accessCount:    make([]uint64, size),
		memorySnapshot: make([]byte, size),
	}

	res.ClearStatistics()

	return res
}

func (l *LinearMemory) TakeSnapshot() {
	copy(l.memorySnapshot, l.memory)
}

func (l *LinearMemory) RestoreSnapshot() {
	copy(l.memory, l.memorySnapshot)
}

func (l *LinearMemory) ClearStatistics() {
	for count := 0; count < len(l.memory); count++ {
		l.accessCount[count] = 0
	}
}

func (l *LinearMemory) Load(address uint16) uint8 {
	l.accessCount[address]++
	return l.memory[address]
}

func (l *LinearMemory) Store(address uint16, b uint8) {
	l.accessCount[address]++
	l.memory[address] = b
}

func (l *LinearMemory) GetStatistics(address uint16) uint64 {
	return l.accessCount[address]
}

func (l *LinearMemory) LoadLarge(address uint32) uint8 {
	return l.Load((uint16)(address & 0xFFFF))
}

func (l *LinearMemory) StoreLarge(address uint32, b uint8) {
	l.Store((uint16)(address&0xFFFF), b)
}

func (l *LinearMemory) GetStatisticsLarge(address uint32) uint64 {
	return l.GetStatistics((uint16)(address & 0xFFFF))
}

func (l *LinearMemory) ToLargeMemory() LargeMemory {
	return l
}
