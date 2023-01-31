package memory

type LinearMemory struct {
	memory      []byte
	accessCount []uint64
}

func NewLinearMemory(size uint16) *LinearMemory {
	res := &LinearMemory{
		memory:      make([]byte, size),
		accessCount: make([]uint64, size),
	}

	res.ClearPerformanceData()

	return res
}

func (l *LinearMemory) ClearPerformanceData() {
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

func (l *LinearMemory) GetStatistics() []BankAccessStatistics {
	memLen := uint16(len(l.memory))

	accessData := BankAccessStatistics{}
	accessData.Address = 0x0000
	accessData.Length = memLen
	accessData.AccessStatistics = make([]uint64, memLen)

	for count := 0; count < len(l.memory); count++ {
		accessData.AccessStatistics[count] = l.accessCount[count]
	}

	return []BankAccessStatistics{accessData}
}
