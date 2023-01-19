package memory

import (
	"encoding/hex"
	"fmt"
)

type Memory interface {
	Load(address uint16) uint8
	Store(address uint16, b uint8)
	GetStatistics()
}

func Dump(m Memory, start uint16, end uint16) {
	temp := make([]byte, end-start+1)
	var index uint16 = 0

	for count := start; count <= end; count++ {
		temp[index] = m.Load(count)
		index++
	}

	fmt.Print(hex.Dump(temp))
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
