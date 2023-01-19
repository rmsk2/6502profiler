package memory

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, b uint8)
	GetStatistics()
}

type LinearMemory struct {
	memory []byte
}

func NewLinearMemory(size uint16) *LinearMemory {
	return &LinearMemory{
		memory: make([]byte, size),
	}
}

func (l *LinearMemory) Load(address uint16) uint8 {
	return l.memory[address]
}

func (l *LinearMemory) Store(address uint16, b uint8) {
	l.memory[address] = b
}

func (l *LinearMemory) GetStatistics() {

}
